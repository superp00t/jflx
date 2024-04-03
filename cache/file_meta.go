package cache

import (
	"encoding/json"
	"os"
	"time"
)

type cache_file_metadata struct {
	Realpath  string            `json:"realpath"`
	Info      cache_file_info   `json:"info"`
	Epoch     time.Time         `json:"epoch"`
	Directory []cache_file_info `json:"directory,omitempty"`
}

func (file *cache_file) read_meta_file() (err error) {
	meta_file_path := file.subfile_path(file.hash + ".met")

	var meta_file []byte
	meta_file, err = os.ReadFile(meta_file_path)
	if err != nil {
		return
	}

	err = json.Unmarshal(meta_file, &file.meta)
	return
}

func (file *cache_file) write_meta_file() (err error) {
	meta_file_path := file.subfile_path(file.hash + ".met")
	var meta_file []byte
	meta_file, err = json.Marshal(&file.meta)
	if err != nil {
		return
	}
	return os.WriteFile(meta_file_path, meta_file, 0700)
}
