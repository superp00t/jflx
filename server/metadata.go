package server

import (
	"bytes"
	"fmt"
	"image"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/superp00t/jflx/conf"
	"github.com/superp00t/jflx/meta"
	"github.com/superp00t/jflx/meta/jflxmeta"

	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
)

var movieRgx = regexp.MustCompile("^(.*) \\((\\d+)\\)$")

func (s *Server) DownloadJPEG(url string) ([]byte, error) {
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

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		panic(err)
	}

	return buf.Bytes(), err
}

func (s *Server) ScrapeMovieDir(path string, info fs.FileInfo) error {
	// Interpret movie directory name
	matches := movieRgx.FindStringSubmatch(info.Name())

	// Find the year the movie was released
	i, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return err
	}

	// Read existing movie metadata
	direc := jflxmeta.Read(filepath.Join(path, "movie.jflxmeta"))

	q := &meta.MovieQuestion{
		Name: matches[1],
		Year: int(i),
	}

	if direc != nil {
		// We are not overwriting curated content
		if direc.DoNotOverwrite {
			return nil
		}

		// If we are OK with overwriting, overwrite with a confirmed ID
		if direc.ForceMovieID != 0 {
			q.HintID = int(direc.ForceMovieID)
		}
	}

	movi, err := s.scraper.AskMovie(q)
	if err != nil {
		log.Println(err)
		return err
	}

	if movi == nil {
		return nil
	}

	log.Println(info.Name())

	nfoPath := filepath.Join(path, "movie.nfo")
	log.Println("writing to", nfoPath)

	if err := movi.Movie.Write(nfoPath); err != nil {
		return err
	}

	os.Remove(filepath.Join(path, "banner.jpg"))

	if movi.PosterURL != "" {
		img, err := s.DownloadJPEG(movi.PosterURL)
		if err == nil {
			ioutil.WriteFile(
				filepath.Join(path, "poster.jpg"),
				img,
				0700)
		}
	}

	if movi.BackdropURL != "" {
		img, err := s.DownloadJPEG(movi.BackdropURL)
		if err == nil {
			ioutil.WriteFile(
				filepath.Join(path, "landscape.jpg"),
				img,
				0700)
		}
	}

	if movi.LogoURL != "" {
		img, err := s.DownloadJPEG(movi.LogoURL)
		if err == nil {
			ioutil.WriteFile(
				filepath.Join(path, "clearlogo.jpg"),
				img,
				0700)
		}
	}

	return nil
}

func (s *Server) ScrapeMovieSource(source string) {
	log.Println("Scraping source", source)

	if err := filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			if movieRgx.MatchString(info.Name()) {
				if err := s.ScrapeMovieDir(path, info); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		log.Println("ScrapeMovieSource: ", err)
	}
}

func (s *Server) perform_scrape() error {
	if s.scrape_in_progress() {
		return fmt.Errorf("already scraping")
	}

	s.scraper_status.Store(true)

	var err error
	s.scraper, err = meta.NewTMDBScraper(s.Conf.TMDBScrapeKey, "")
	if err != nil {
		s.scraper = nil
		s.scraper_status.Store(false)
		return err
	}

	for _, vol := range s.Volumes {
		if vol.Conf.Kinds == conf.Movie {
			for _, src := range vol.Conf.Sources {
				s.ScrapeMovieSource(src)
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
