package server

import (
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/superp00t/jflx/conf"
)

type Volume struct {
	Conf    *conf.Volume
	Handler http.Handler
}

type volumeDir struct {
	v       *Volume
	dirname string
}

func (vd *volumeDir) Close() error {
	return nil
}

func (vd *volumeDir) Read(b []byte) (int, error) {
	return 0, errors.New("JFLX/server: cannot read from directory file")
}

func (vd *volumeDir) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("JFLX/server: cannot seek from directory file")
}

func (vd *volumeDir) Readdir(count int) ([]fs.FileInfo, error) {
	log.Println("read dir", vd.dirname)

	var ls []fs.FileInfo

	for _, d := range vd.v.Conf.Sources {
		dir := string(d)
		if dir == "" {
			dir = "."
		}
		fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+vd.dirname)))
		dirents, err := os.ReadDir(fullName)
		if err == nil {
			for _, ent := range dirents {
				info, err := ent.Info()
				if err == nil {
					ls = append(ls, info)
				}
			}
		}
	}

	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Name() <= ls[j].Name()
	})

	if count > 0 {
		ls = ls[:count]
	}

	return ls, nil
}

func (vd *volumeDir) Stat() (fs.FileInfo, error) {
	for _, d := range vd.v.Conf.Sources {
		dir := string(d)
		if dir == "" {
			dir = "."
		}
		fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+vd.dirname)))

		f, err := os.Open(fullName)
		if err != nil {
			continue
		}
		st, err := f.Stat()
		if err != nil {
			f.Close()
			continue
		}
		f.Close()
		return st, nil
	}

	return nil, errors.New("JFLX/server: volumeDir stat mismatch")
}

func (v *Volume) openDirectoryFile(name string) (http.File, error) {
	return &volumeDir{v, name}, nil
}

func (v *Volume) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("JFLX/server: invalid character in file path")
	}

	for _, d := range v.Conf.Sources {
		dir := string(d)
		if dir == "" {
			dir = "."
		}
		fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
		f, err := os.Open(fullName)
		if err != nil {
			continue
		}
		if stat, err := f.Stat(); err == nil {
			if stat.IsDir() {
				f.Close()
				return v.openDirectoryFile(name)
			}
		}
		return f, nil
	}

	return nil, os.ErrNotExist
}
