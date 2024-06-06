package cache

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	max_cached_file_size int64 = 0xFFFFFFFFF
	file_part_size       int64 = 0xFFFFF

	free_time_interval = 16 * time.Minute
)

type Config struct {
	Directory       string
	MaxAge          time.Duration
	MaxDirectoryAge time.Duration
	MaxSize         int64
}

type Server struct {
	guard_files    sync.Mutex
	config         *Config
	source         http.FileSystem
	used_bytes     atomic.Int64
	last_free_time time.Time
}

func NewServer(config *Config, source http.FileSystem) (s *Server, err error) {
	s = new(Server)
	s.config = config
	s.source = source

	if s.config.Directory == "" {
		err = fmt.Errorf("cache directory must have a value")
		return
	}

	if err = os.MkdirAll(s.config.Directory, 0700); err != nil {
		return
	}

	if err = s.compute_used_bytes(); err != nil {
		return
	}

	// TODO: cleanup all dirty files
	return
}
