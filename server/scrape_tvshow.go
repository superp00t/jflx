package server

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/superp00t/jflx/meta"
	"github.com/superp00t/jflx/meta/jflxmeta"
	"github.com/superp00t/jflx/meta/nfo"
)

func (s *Server) scrape_tvshow_season_episode(tvshow_ID int, season int, season_directory string, episode_number_in_season int, episode_filename string, episode_title string) error {
	var q meta.EpisodeQuestion
	q.Episode = episode_number_in_season
	q.Season = season
	q.ShowID = tvshow_ID
	q.HintEpisodeTitle = episode_title

	extension := filepath.Ext(episode_filename)
	episode_filename_without_extension := episode_filename[0 : len(episode_filename)-len(extension)]

	nfo_path := filepath.Join(season_directory, episode_filename_without_extension+".nfo")
	log.Println("writing to", nfo_path)

	episode, err := s.scraper.AskTvshowEpisode(&q)
	if err != nil {
		log.Println(err)
		return nil
	}
	if err := nfo.WriteTvshowEpisode(nfo_path, &episode.TvshowEpisode); err != nil {
		return err
	}

	os.Remove(filepath.Join(season_directory, episode_filename_without_extension+"-banner.jpg"))

	if episode.PosterURL != "" {
		if err = s.update_jpeg_file(poster_artwork, episode.PosterURL, filepath.Join(season_directory, episode_filename_without_extension+"-poster.jpg")); err != nil {
			return err
		}
	}

	if episode.BackdropURL != "" {
		if err = s.update_jpeg_file(backdrop_artwork, episode.BackdropURL, filepath.Join(season_directory, episode_filename_without_extension+"-landscape.jpg")); err != nil {
			return err
		}
	}

	if episode.LogoURL != "" {
		if err = s.update_jpeg_file(logo_artwork, episode.LogoURL, filepath.Join(season_directory, episode_filename_without_extension+"-clearlogo.jpg")); err != nil {
			return err
		}
	}

	if episode.ThumbURL != "" {
		if err = s.update_jpeg_file(thumb_artwork, episode.ThumbURL, filepath.Join(season_directory, episode_filename_without_extension+"-thumb.jpg")); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) scrape_tvshow_season_directory(tvshow_ID int, season int, path string) error {
	season_list, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, item := range season_list {
		if !item.IsDir() {
			if episodic_regex.MatchString(item.Name()) {
				matches := episodic_regex.FindStringSubmatch(item.Name())

				episode_number, err := strconv.ParseInt(matches[2], 10, 64)
				if err != nil {
					log.Println("matches are", spew.Sdump(matches))
					return err
				}

				episode_title := matches[3]
				if err := s.scrape_tvshow_season_episode(tvshow_ID, season, path, int(episode_number), item.Name(), episode_title); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *Server) scrape_tvshow_directory(path string, name string) error {
	// Interpret tvshow directory name
	matches := name_and_year_regex.FindStringSubmatch(name)

	// Find the year the tvshow was released
	i, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return err
	}

	// Skip already scraped tvshows
	if _, err := os.Stat(filepath.Join(path, "tvshow.nfo")); err == nil {
		return nil
	}

	// Read existing tvshow metadata
	tvshow_meta, _ := jflxmeta.ReadTvshow(filepath.Join(path, "tvshow.jflxmeta"))

	//
	nfo_path := filepath.Join(path, "tvshow.nfo")

	// Skip if NFO has same ID as meta file
	if _, err := os.Stat(nfo_path); err == nil {
		var nfo_data nfo.Tvshow
		err := nfo.ReadTvshow(nfo_path, &nfo_data)
		if err != nil {
			return err
		}

		if tvshow_meta != nil {
			if nfo.GetDefault(nfo_data.Uniqueids).Text == strconv.FormatInt(tvshow_meta.ForceTvshowID, 10) {
				return nil
			}
		}
	}

	q := &meta.ShowQuestion{
		Name: matches[1],
		Year: int(i),
	}

	show_id := int(0)

	if tvshow_meta != nil {
		// We are not overwriting curated content
		if tvshow_meta.DoNotOverwrite {
			return nil
		}

		// If we are OK with overwriting, overwrite with a confirmed ID
		if tvshow_meta.ForceTvshowID != 0 {
			q.HintID = int(tvshow_meta.ForceTvshowID)
		}
	}

	show, err := s.scraper.AskTvshow(q)
	if err != nil {
		log.Println(err)
		return err
	}

	if show == nil {
		return nil
	}

	show_id, err = strconv.Atoi(nfo.GetDefault(show.Tvshow.Uniqueids).Text)
	if err != nil {
		panic(err)
	}

	if err := nfo.WriteTvshow(nfo_path, &show.Tvshow); err != nil {
		return err
	}

	os.Remove(filepath.Join(path, "banner.jpg"))

	if show.PosterURL != "" {
		if err := s.update_jpeg_file(poster_artwork, show.PosterURL, filepath.Join(path, "poster.jpg")); err != nil {
			return err
		}
	}

	if show.BackdropURL != "" {
		if err := s.update_jpeg_file(backdrop_artwork, show.BackdropURL, filepath.Join(path, "landscape.jpg")); err != nil {
			return err
		}
	}

	if show.LogoURL != "" {
		if err := s.update_jpeg_file(logo_artwork, show.LogoURL, filepath.Join(path, "clearlogo.jpg")); err != nil {
			return err
		}
	}

	// scan seasons
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, ent := range dir {
		if ent.IsDir() {
			// Is this a valid season name? 00-99
			if season_folder_regex.MatchString(ent.Name()) {

				season, err := strconv.ParseInt(ent.Name(), 10, 64)
				if err == nil {
					if err = s.scrape_tvshow_season_directory(show_id, int(season), filepath.Join(path, ent.Name())); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (s *Server) scrape_tvshow_source(source string) {
	log.Println("Scraping source", source)

	source_list, err := os.ReadDir(source)
	if err != nil {
		log.Println("scrape_tvshow_source: ReadDir(): ", err)
		return
	}

	for _, info := range source_list {
		if info.IsDir() {
			if name_and_year_regex.MatchString(info.Name()) {
				if err := s.scrape_tvshow_directory(filepath.Join(source, info.Name()), info.Name()); err != nil {
					log.Println("scrape_tvshow_directory: ", err)
					return
				}
			}
		}
	}
}
