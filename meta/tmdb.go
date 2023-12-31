package meta

import (
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
)

type TMDBScraper struct {
	Client *tmdb.Client
}

func NewTMDBScraper(key, userAgent string) (t *TMDBScraper, err error) {
	t = new(TMDBScraper)
	t.Client, err = tmdb.Init(key)
	return
}

func (t *TMDBScraper) AskMovieID(q *MovieQuestion) (id int, err error) {
	sm, err := t.Client.GetSearchMovies(q.Name, map[string]string{
		"year": strconv.FormatInt(int64(q.Year), 10),
	})
	if err != nil {
		return 0, err
	}
	if len(sm.Results) == 0 {
		return 0, nil
	}
	return int(sm.Results[0].ID), nil
}

func (t *TMDBScraper) AskMovie(q *MovieQuestion) (*MovieAnswer, error) {
	ma := new(MovieAnswer)

	var id int

	if q.HintID != 0 {
		id = q.HintID
	} else {
		var err error
		id, err = t.AskMovieID(q)
		if err != nil {
			return nil, err
		}
	}

	if id == 0 {
		return nil, nil
	}

	details, err := t.Client.GetMovieDetails(id, nil)
	if err != nil {
		return nil, err
	}

	ma.Title = details.Title
	ma.Originaltitle = details.OriginalTitle
	ma.Plot = details.Overview
	ma.Tagline = details.Tagline
	ma.Premiered = details.ReleaseDate

	if details.PosterPath != "" {
		ma.PosterURL = tmdb.GetImageURL(details.PosterPath, tmdb.Original)
	}

	if details.BackdropPath != "" {
		ma.BackdropURL = tmdb.GetImageURL(details.BackdropPath, tmdb.Original)
	}

	if details.MovieImagesAppend != nil && details.MovieImagesAppend.Images != nil {
		if len(details.Images.Logos) > 0 {
			ma.LogoURL = tmdb.GetImageURL(details.Images.Logos[0].FilePath, tmdb.Original)
		}
	}

	return ma, nil
}
