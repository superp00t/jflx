package server

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"net/http"
	"regexp"
	"time"

	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
	"github.com/superp00t/jflx/conf"
	"github.com/superp00t/jflx/meta"
)

var (
	name_and_year_regex = regexp.MustCompile(`^(.*) \((\d+)\)$`)
	season_folder_regex = regexp.MustCompile("[0-9]+")
	episodic_regex      = regexp.MustCompile(`^S([0-9]{2})E([0-9]{2})\s?(.*)\.(mkv|mp4|avi|mov|mpeg|ts|webm)$`)
)

func (s *Server) download_jpeg_url(at art_type, url string) ([]byte, error) {
	cl := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	if at.constrained() {
		img = resize.Resize(uint(at.Width), uint(at.Height), img, resize.Lanczos3)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{
		Quality: 75,
	}); err != nil {
		panic(err)
	}

	return buf.Bytes(), err
}

func (s *Server) perform_scrape() error {
	if s.scrape_in_progress() {
		return fmt.Errorf("already scraping")
	}

	s.scraper_status.Store(true)

	var err error
	s.scraper, err = meta.NewTMDBScraper(s.config.TMDBScrapeKey, "")
	if err != nil {
		s.scraper = nil
		s.scraper_status.Store(false)
		return err
	}

	// Scrape movie sources
	for _, volume_config := range s.config.Volumes {
		if volume_config.Kinds == conf.Movie {
			for _, src := range volume_config.Sources {
				s.scrape_movie_source(src)
			}
		}
	}

	for _, volume_config := range s.config.Volumes {
		if volume_config.Kinds == conf.TvShow {
			for _, src := range volume_config.Sources {
				s.scrape_tvshow_source(src)
			}
		}
	}

	s.scraper = nil
	s.scraper_status.Store(false)

	return nil
}

func (s *Server) perform_scrape_and_log() {
	err := s.perform_scrape()
	if err == nil {
		log.Println("scrape successful!")
	} else {
		log.Println("scrape failed:", err)
	}
}

func (s *Server) scrape_in_progress() bool {
	return s.scraper_status.Load()
}
