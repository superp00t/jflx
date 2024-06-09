package meta

import (
	tmdb "github.com/cyruzin/golang-tmdb"
)

type TMDBScraper struct {
	cached_tv_details map[int]*tmdb.TVDetails
	Client            *tmdb.Client
}

func NewTMDBScraper(key, userAgent string) (t *TMDBScraper, err error) {
	t = new(TMDBScraper)
	t.cached_tv_details = make(map[int]*tmdb.TVDetails)
	t.Client, err = tmdb.Init(key)
	return
}
