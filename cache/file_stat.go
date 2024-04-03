package cache

import "io/fs"

func (file *cache_file) Stat() (info fs.FileInfo, err error) {
	return &file.meta.Info, nil
}
