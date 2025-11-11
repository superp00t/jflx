package scraper

import (
	"github.com/superp00t/jflx/config"
	"github.com/superp00t/jflx/meta"
)

type ScrapeLibraryParams struct {
	TMDB_API_Key string
	Volumes      []config.VolumeConfig
}

func ScrapeLibrary(params *ScrapeLibraryParams) (err error) {
	var s Scraper
	s.meta_db, err = meta.NewDatabaseTMDB(params.TMDB_API_Key, "")
	if err != nil {
		return
	}
	err = s.perform_library_scrape(params.Volumes)
	return
}
