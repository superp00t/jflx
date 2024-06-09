package jflxmeta

import (
	"encoding/json"
	"os"
)

type Movie struct {
	DoNotOverwrite bool  `json:"do_not_overwrite"`
	ForceMovieID   int64 `json:"force_movie_ID"`
}

func ReadMovie(path string) (m *Movie, err error) {
	var b []byte
	b, err = os.ReadFile(path)
	if err != nil {
		return
	}

	m = new(Movie)
	if err := json.Unmarshal(b, m); err != nil {
		panic(err)
	}
	return m, nil
}
