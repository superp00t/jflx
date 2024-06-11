package meta

import (
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/superp00t/jflx/meta/nfo"
)

func (t *TMDBScraper) ask_movie_ID(q *ShowQuestion) (id int, err error) {
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

func (t *TMDBScraper) AskMovie(q *ShowQuestion) (*MovieAnswer, error) {
	ma := new(MovieAnswer)

	var id int

	if q.HintID != 0 {
		id = q.HintID
	} else {
		var err error
		id, err = t.ask_movie_ID(q)
		if err != nil {
			return nil, err
		}
	}

	if id == 0 {
		return nil, nil
	}

	details, err := t.Client.GetMovieDetails(id, map[string]string{
		"append_to_response": "credits",
	})
	if err != nil {
		return nil, err
	}

	ma.Uniqueids = []nfo.ID{
		{
			Type:    "tmdb",
			Default: "true",
			Text:    strconv.FormatInt(int64(id), 10),
		},
	}

	for _, credit := range details.Credits.Crew {
		switch credit.Job {
		case "Director":
			ma.Directors = append(ma.Directors, credit.Name)
		}
	}

	ma.Title = details.Title
	ma.Originaltitle = details.OriginalTitle
	ma.Plot = details.Overview
	ma.Tagline = details.Tagline
	ma.Premiered = details.ReleaseDate
	ma.Tagline = details.Tagline

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
