package cache

import (
	"errors"
	"io/fs"
)

func (file *cache_file) Readdir(count int) (directory []fs.FileInfo, err error) {
	if !file.meta.Info.IsDir() {
		err = errors.New("cache_file.Readdir: cannot readdir on a normal file")
		return
	}
	if count == -1 || count > len(file.meta.Directory) {
		directory = make([]fs.FileInfo, len(file.meta.Directory))
	} else {
		directory = make([]fs.FileInfo, count)
	}
	for i := range directory {
		directory[i] = &file.meta.Directory[i]
	}
	return
}
