package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"time"
)

// file structure:
// - ab
//  - cd
//   - abcdef... .lat          -- the last accessed time in JSON
//   - abcdef... .met          -- JSON representation for file/directory
//   - abcdef... -00000000.dat -- data segment #1
//   - abcdef... -00000001.dat -- data segment #2
//

type cache_file struct {
	server           *Server
	meta             cache_file_metadata
	last_access_time time.Time
	// the path hash hex prefix in the cache directory e.g. (ab/cd)
	prefix string
	// the full hash of the filename e.g. (abcdef...)
	hash     string
	realpath string
	pointer  int64
}

func (file *cache_file) subfile_path(path string) string {
	return filepath.Join(file.server.config.Directory, file.prefix, path)
}

func get_cached_filepath(path string) (prefix string, cached_path string) {
	h := sha256.New()
	h.Write([]byte(path))
	digest := h.Sum(nil)
	hexed := hex.EncodeToString(digest)

	path1 := hexed[0:2]
	path2 := hexed[2:4]

	prefix = fmt.Sprintf("%s/%s", path1, path2)
	cached_path = hexed
	return
}
