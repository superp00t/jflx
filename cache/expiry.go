package cache

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// remove a cached file from the cache_directory
func (s *Server) free_cached_file(prefix, hash string) (freed_bytes int64, err error) {
	var directory_entities []os.DirEntry

	prefixed_directory := filepath.Join(s.config.Directory, prefix)

	directory_entities, err = os.ReadDir(prefixed_directory)
	if err != nil {
		// if directory doesn't exist, okay
		err = nil
		return
	}

	for _, entity := range directory_entities {
		name := entity.Name()
		if strings.HasPrefix(name, hash) {
			var size int64
			var fi fs.FileInfo
			fi, err = entity.Info()
			if err == nil {
				size = fi.Size()
			}
			if err = os.Remove(filepath.Join(prefixed_directory, name)); err != nil {
				return
			}
			freed_bytes += size
		}
	}

	log.Println("freed cache file", hash)

	s.add_used_bytes(-freed_bytes)

	return
}

type cache_age struct {
	hash             string
	last_access_time time.Time
}

// checks if the cache needs to free space, and deletes the least recently used cache files
func (s *Server) free_oldest_files() (err error) {
	var cache_lru []cache_age

	// walk entire cache, scanning both size of used bytes, and sorting the least recently accessed files first

	if err = filepath.Walk(s.config.Directory, func(path string, info fs.FileInfo, err error) error {
		if err == nil {
			fmt.Println("attempting free", path)
			if strings.HasSuffix(path, ".lat") {
				if len(path) >= (4 + 64) {
					var age cache_age
					start_of_hash := len(path) - (4 + 64)
					age.hash = path[start_of_hash : start_of_hash+64]
					age.last_access_time, err = read_last_access_time(path)
					if err == nil {
						cache_lru = append(cache_lru, age)
					}
				}
			}
		}
		return nil
	}); err != nil {
		return
	}

	sort.Slice(cache_lru, func(i, j int) bool {
		return cache_lru[i].last_access_time.Before(cache_lru[j].last_access_time)
	})

	// Free least recently used files first
	for _, age := range cache_lru {
		// If we've met our space limitations, stop now
		if !s.is_overweight() {
			break
		}

		var freed_bytes int64
		prefix := fmt.Sprintf("%s/%s", age.hash[0:2], age.hash[2:4])
		freed_bytes, err = s.free_cached_file(prefix, age.hash)
		if err != nil {
			return
		}
		log.Println("freed", freed_bytes, "bytes")
	}

	return
}

func (s *Server) is_cache_locked() bool {
	return s.cache_locked.Load()
}

func (s *Server) lock_cache() bool {
	if !s.is_cache_locked() {
		return false
	}

	success := s.guard_files.TryLock()
	if success {
		s.cache_locked.Store(true)
	}

	return success
}

func (s *Server) unlock_cache() {
	s.guard_files.Unlock()
	s.cache_locked.Store(false)
}

func (s *Server) attempt_free() (err error) {
	if s.is_cache_locked() {
		return nil
	}

	if time.Since(s.last_free_time) > free_time_interval && s.is_overweight() {
		if !s.lock_cache() {
			return
		}

		err = s.free_oldest_files()

		s.last_free_time = time.Now()
		s.unlock_cache()
	}

	return
}

func (s *Server) compute_used_bytes() (err error) {
	used_bytes := int64(0)

	log.Println("computing cache size, please wait....")

	if err = filepath.Walk(s.config.Directory, func(path string, info fs.FileInfo, err error) error {
		if err == nil {
			used_bytes += info.Size()
		}
		return nil
	}); err != nil {
		return
	}

	s.used_bytes.Store(used_bytes)

	log.Printf("Ready! cache size is %d (%f%% of max capacity)", s.used_bytes.Load(), s.weight()*100.0)

	return
}

func (s *Server) add_used_bytes(delta int64) {
	s.used_bytes.Add(delta)
}

func (s *Server) is_overweight() bool {
	return s.used_bytes.Load() > s.config.MaxSize
}

func (s *Server) weight() float32 {
	return float32(s.used_bytes.Load()) / float32(s.config.MaxSize)
}
