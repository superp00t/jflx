package scraper

import (
	"regexp"

	_ "image/jpeg"
	_ "image/png"

	"github.com/superp00t/jflx/config"
)

var (
	name_and_year_regex = regexp.MustCompile(`^(.*) \((\d+)\)$`)
	season_folder_regex = regexp.MustCompile("[0-9]+")
	episodic_regex      = regexp.MustCompile(`^S([0-9]{2})E([0-9]{2})\s?(.*)\.(mkv|mp4|avi|mov|mpeg|ts|webm)$`)
)

func (s *Scraper) perform_library_scrape(volumes []config.VolumeConfig) error {
	// Scrape movie sources
	for _, volume_config := range volumes {
		if volume_config.Kinds == config.Movie {
			for _, src := range volume_config.Sources {
				s.scrape_movie_source(src)
			}
		}
	}

	for _, volume_config := range volumes {
		if volume_config.Kinds == config.TvShow {
			for _, src := range volume_config.Sources {
				s.scrape_tvshow_source(src)
			}
		}
	}

	return nil
}
