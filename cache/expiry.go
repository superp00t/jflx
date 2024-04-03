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

	return
}

type cache_age struct {
	hash             string
	last_access_time time.Time
}

// checks if the cache needs to free space, and deletes the least recently used cache files
func (s *Server) free_oldest_files() (err error) {
	var used_bytes int64
	var cache_lru []cache_age

	// walk entire cache, scanning both size of used bytes, and sorting the least recently accessed files first

	if err = filepath.Walk(s.config.Directory, func(path string, info fs.FileInfo, err error) error {
		if err == nil {
			used_bytes += info.Size()
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
		overweight := used_bytes >= s.config.MaxSize
		if !overweight {
			break
		}

		var freed_bytes int64
		prefix := fmt.Sprintf("%s/%s", age.hash[0:2], age.hash[2:4])
		freed_bytes, err = s.free_cached_file(prefix, age.hash)
		if err != nil {
			return
		}
		used_bytes -= freed_bytes
		log.Println("freed", freed_bytes, "bytes")
	}

	return
}

func (s *Server) attempt_free() (err error) {
	if time.Since(s.last_free_time) > free_time_interval {
		if !s.guard_files.TryLock() {
			return
		}

		err = s.free_oldest_files()

		s.last_free_time = time.Now()
		s.guard_files.Unlock()
	}

	return
}
