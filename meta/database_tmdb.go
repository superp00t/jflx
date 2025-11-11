package meta

import (
	tmdb "github.com/cyruzin/golang-tmdb"
)

type DatabaseTMDB struct {
	cached_tv_details map[int]*tmdb.TVDetails
	Client            *tmdb.Client
}

func NewDatabaseTMDB(key, userAgent string) (t *DatabaseTMDB, err error) {
	t = new(DatabaseTMDB)
	t.cached_tv_details = make(map[int]*tmdb.TVDetails)
	t.Client, err = tmdb.Init(key)
	return
}
