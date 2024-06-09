package cache

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Create a virtual cached file
func (s *Server) open_cache_file(path string) (file *cache_file, err error) {
	access_time := time.Now()

	// create new file instance
	file = new(cache_file)
	file.server = s
	file.realpath = path
	file.prefix, file.hash = get_cached_filepath(path)

	var (
		source_file        http.File
		fi                 fs.FileInfo
		opened_source_file bool
	)

	// check if meta file exists
	cached_file_exists := false
	valid_meta_file := false
	if err = file.read_meta_file(); err != nil {
		// If not, we need to create one
		cached_file_exists = false
		valid_meta_file = false
	} else {
		cached_file_exists = true

		// A meta file exists, but is it expired?

		// See if it's recent enough to just proceed. Checking Stat can stress out your storage media
		valid_age := s.config.MaxAge
		if file.meta.Info.InfoIsDirectory {
			valid_age = s.config.MaxDirectoryAge
		}

		invalidation_time := file.meta.Epoch.Add(valid_age)

		if access_time.Compare(invalidation_time) >= 0 {
			// log.Println("Cached metadata for ", file.realpath, "may be too old")

			valid_meta_file = false

		} else {
			// log.Println("Cached metadata for ", file.realpath, "is certainly not too old, ", invalidation_time.Sub(access_time), " until invalidation")
			valid_meta_file = true
		}
	}

	// File has been invalidated
	if cached_file_exists && !valid_meta_file {
		s.guard_files.Lock()
		if _, err = s.free_cached_file(file.prefix, file.hash); err != nil {
			s.guard_files.Unlock()
			return
		}
		s.guard_files.Unlock()
	}

	if !valid_meta_file {
		// if file isn't cached, make metadata file now
		if !opened_source_file {
			source_file, err = s.source.Open(path)
			if err != nil {
				return
			}
			fi, err = source_file.Stat()
			if err != nil {
				return
			}
			opened_source_file = true
			defer source_file.Close()
		}

		// create prefix directory if file exists
		prefix_directory := filepath.Join(s.config.Directory, file.prefix)
		if _, err = os.Stat(prefix_directory); err != nil {
			// prefix folder doesn't exist
			err = nil
			if err = os.MkdirAll(prefix_directory, 0700); err != nil {
				return
			}
		}

		file.meta.Realpath = path

		file.meta.Epoch = time.Now()
		file.meta.Info.InfoIsDirectory = fi.IsDir()
		file.meta.Info.InfoName = fi.Name()
		file.meta.Info.InfoModTime = fi.ModTime()
		file.meta.Info.InfoSize = fi.Size()
		file.meta.Info.InfoMode = fi.Mode()

		if file.meta.Info.IsDir() {
			var directory_entities []fs.FileInfo
			directory_entities, err = source_file.Readdir(-1)
			if err != nil {
				return
			}

			file.meta.Directory = make([]cache_file_info, len(directory_entities))
			for i := range file.meta.Directory {
				directory_file_info := &file.meta.Directory[i]
				directory_entity := directory_entities[i]
				directory_file_info.InfoIsDirectory = directory_entity.IsDir()
				directory_file_info.InfoName = directory_entity.Name()
				directory_file_info.InfoModTime = directory_entity.ModTime()
				directory_file_info.InfoSize = directory_entity.Size()
				directory_file_info.InfoMode = directory_entity.Mode()
			}
		}

		if err = file.write_meta_file(); err != nil {
			return
		}
	}

	// mark that we accessed it now
	if err = file.write_last_access_time(time.Now()); err != nil {
		return
	}

	// try to sweep
	if err = s.attempt_free(); err != nil {
		return
	}

	return
}

func (s *Server) Open(path string) (file http.File, err error) {
	avoid_cache := false

	if s.is_cache_locked() {
		// cache is being freed.
		avoid_cache = true
	} else if s.is_overweight() {
		avoid_cache = true
		go s.attempt_free()
	}

	if avoid_cache {
		return s.source.Open(path)
	}

	file, err = s.open_cache_file(path)
	if err != nil {
		return nil, err
	}

	return
}
