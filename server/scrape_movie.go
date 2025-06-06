package server

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/superp00t/jflx/meta"
	"github.com/superp00t/jflx/meta/jflxmeta"
	"github.com/superp00t/jflx/meta/nfo"
)

func (s *Server) scrape_movie_directory(path string, info fs.FileInfo) error {
	// Interpret movie directory name
	matches := name_and_year_regex.FindStringSubmatch(info.Name())

	// Find the year the movie was released
	i, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return err
	}

	// Read existing movie metadata
	movie_meta, _ := jflxmeta.ReadMovie(filepath.Join(path, "movie.jflxmeta"))

	//
	nfo_path := filepath.Join(path, "movie.nfo")

	id_changed := false

	// Skip if NFO has same ID as meta file
	if _, err := os.Stat(nfo_path); err == nil {
		// If NFO file exists
		var nfo_data nfo.Movie
		err := nfo.ReadMovie(nfo_path, &nfo_data)
		if err == nil {
			// If NFO file has an ID
			default_id := nfo.GetDefault(nfo_data.Uniqueids)

			if default_id != nil {
				// If movie has a JFLXMETA file
				if movie_meta != nil {
					// if movie has the same ID as the forced one, nothing to do
					if default_id.Text == strconv.FormatInt(movie_meta.ForceMovieID, 10) {
						return nil
					}
					id_changed = true
					// otherwise, refresh
				} else {
					// if not, skip. (we already have NFO)
					return nil
				}
			}
		}
	}

	q := &meta.ShowQuestion{
		Name: matches[1],
		Year: int(i),
	}

	if movie_meta != nil {
		// We are not overwriting curated content
		if movie_meta.DoNotOverwrite {
			return nil
		}

		// If we are OK with overwriting, overwrite with a confirmed ID
		if movie_meta.ForceMovieID != 0 {
			q.HintID = int(movie_meta.ForceMovieID)
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

	if err := nfo.WriteMovie(nfo_path, &movi.Movie); err != nil {
		return err
	}

	if id_changed {
		os.Remove(filepath.Join(path, "banner.jpg"))
		os.Remove(filepath.Join(path, "poster.jpg"))
		os.Remove(filepath.Join(path, "landscape.jpg"))
		os.Remove(filepath.Join(path, "clearlogo.jpg"))
	}

	if movi.PosterURL != "" {
		if err := s.update_jpeg_file(poster_artwork, movi.PosterURL, filepath.Join(path, "poster.jpg")); err != nil {
			return err
		}
	}

	if movi.BackdropURL != "" {
		if err := s.update_jpeg_file(backdrop_artwork, movi.BackdropURL, filepath.Join(path, "landscape.jpg")); err != nil {
			return err
		}
	}

	if movi.LogoURL != "" {
		if err := s.update_jpeg_file(logo_artwork, movi.LogoURL, filepath.Join(path, "clearlogo.jpg")); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) scrape_movie_source(source string) {
	log.Println("Scraping source", source)

	if err := filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			if name_and_year_regex.MatchString(info.Name()) {
				if err := s.scrape_movie_directory(path, info); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		log.Println("ScrapeMovieSource: ", err)
	}
}
