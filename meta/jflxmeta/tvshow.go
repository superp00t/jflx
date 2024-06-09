package jflxmeta

import (
	"encoding/json"
	"os"
)

type Tvshow struct {
	DoNotOverwrite bool  `json:"do_not_overwrite"`
	AlwaysRefresh  bool  `json:"always_refresh"`
	ForceTvshowID  int64 `json:"force_tvshow_ID"`
}

func ReadTvshow(path string) (m *Tvshow, err error) {
	var b []byte
	b, err = os.ReadFile(path)
	if err != nil {
		return
	}

	m = new(Tvshow)
	if err := json.Unmarshal(b, m); err != nil {
		panic(err)
	}
	return m, nil
}
