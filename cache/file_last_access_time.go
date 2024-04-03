package cache

import (
	"encoding/json"
	"os"
	"time"
)

func read_last_access_time(path string) (t time.Time, err error) {
	var data []byte
	data, err = os.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &t)
	return
}

func write_last_access_time(path string, t time.Time) (err error) {
	var data []byte
	data, err = json.Marshal(&t)
	if err != nil {
		return
	}
	return os.WriteFile(path, data, 0700)
}

func (file *cache_file) read_last_access_time() (err error) {
	file.last_access_time, err = read_last_access_time(file.subfile_path(file.hash + ".lat"))
	return
}

func (file *cache_file) write_last_access_time(t time.Time) (err error) {
	err = write_last_access_time(file.subfile_path(file.hash+".lat"), t)
	return
}
