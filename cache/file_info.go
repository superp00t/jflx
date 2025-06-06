package cache

import (
	"io/fs"
	"time"
)

type cache_file_info struct {
	InfoName        string      `json:"name"`
	InfoMode        fs.FileMode `json:"mode"`
	InfoSize        int64       `json:"size"`
	InfoModTime     time.Time   `json:"mod_time"`
	InfoIsDirectory bool        `json:"is_dir"`
}

func (info *cache_file_info) Name() string {
	return info.InfoName
}

func (info *cache_file_info) Mode() fs.FileMode {
	return info.InfoMode
}

func (info *cache_file_info) Size() int64 {
	return info.InfoSize
}

func (info *cache_file_info) IsDir() bool {
	return info.InfoIsDirectory
}

func (info *cache_file_info) ModTime() time.Time {
	return info.InfoModTime
}

func (info *cache_file_info) Sys() any {
	return nil
}
