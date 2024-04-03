package cache

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Create a virtual cached file
func (s *Server) open_cache_file(path string) (file *cache_file, err error) {
	file = new(cache_file)
	file.server = s
	file.realpath = path
	file.prefix, file.hash = get_cached_filepath(path)

	cached_file_exists := false
	valid_meta_file := true
	if err = file.read_meta_file(); err != nil {
		cached_file_exists = true
		valid_meta_file = false
	} else {
		valid_age := s.config.MaxAge
		if file.meta.Info.InfoIsDirectory {
			valid_age = s.config.MaxDirectoryAge
		}

		if time.Since(file.meta.Epoch) > valid_age {
			valid_meta_file = false
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
		var source_file http.File
		var fi fs.FileInfo
		source_file, err = s.source.Open(path)
		if err != nil {
			return
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

		fi, err = source_file.Stat()
		if err != nil {
			return
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
	file, err = s.open_cache_file(path)
	if err != nil {
		log.Println("failed to open path", path, "err", err)
		// fallback to source file
		// return s.source.Open(path)
		return nil, err
	}

	return
}
