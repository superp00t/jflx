package cache

import (
	"errors"
	"io"
)

func (file *cache_file) Seek(offset int64, whence int) (int64, error) {
	if file.meta.Info.IsDir() {
		return 0, errors.New("cache_file.Seek: cannot seek on directory file")
	}
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = file.pointer + offset
	case io.SeekEnd:
		abs = file.meta.Info.Size() + offset
	default:
		return 0, errors.New("cache_file.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("cache_file.Seek: negative position")
	}
	if abs >= file.meta.Info.Size() {
		return 0, errors.New("cache_file.Seek: cannot seek past end of file")
	}
	file.pointer = abs
	return abs, nil
}
