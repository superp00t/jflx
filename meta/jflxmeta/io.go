package jflxmeta

import (
	"encoding/json"
	"io/ioutil"
)

type Data struct {
	DoNotOverwrite bool  `json:"do_not_overwrite"`
	ForceMovieID   int64 `json:"force_movie_ID"`
}

func Read(path string) *Data {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	d := new(Data)
	if err := json.Unmarshal(b, d); err != nil {
		panic(err)
	}
	return d
}
