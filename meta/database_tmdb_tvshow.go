package meta

import (
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/superp00t/jflx/meta/nfo"
)

func (t *DatabaseTMDB) ask_tv_show_ID(q *ShowQuestion) (id int, err error) {
	sm, err := t.Client.GetSearchTVShow(q.Name, map[string]string{
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

func (t *DatabaseTMDB) get_tv_details(id int) (details *tmdb.TVDetails, err error) {
	var ok bool
	details, ok = t.cached_tv_details[id]
	if ok {
		return
	}

	details, err = t.Client.GetTVDetails(id, nil)
	if err != nil {
		return
	}

	t.cached_tv_details[id] = details
	return
}

func (t *DatabaseTMDB) AskTvshow(q *ShowQuestion) (tvshow *TvshowAnswer, err error) {
	tvshow = new(TvshowAnswer)

	var id int

	if q.HintID != 0 {
		id = q.HintID
	} else {
		var err error
		id, err = t.ask_tv_show_ID(q)
		if err != nil {
			return nil, err
		}
	}

	if id == 0 {
		return nil, nil
	}

	details, err := t.Client.GetTVDetails(id, map[string]string{
		// "append_to_response": "credits",
	})
	if err != nil {
		return nil, err
	}

	tvshow.Uniqueids = []nfo.ID{
		{
			Type:    "tmdb",
			Default: "true",
			Text:    strconv.FormatInt(int64(id), 10),
		},
	}
	tvshow.Title = details.Name
	tvshow.Originaltitle = details.OriginalName
	tvshow.Plot = details.Overview
	tvshow.Tagline = details.Tagline
	tvshow.Premiered = details.FirstAirDate
	if len(details.FirstAirDate) > 4 {
		tvshow.Year = details.FirstAirDate[0:4]
	}
	tvshow.Episode = strconv.FormatInt(int64(details.NumberOfEpisodes), 10)
	tvshow.Season = strconv.FormatInt(int64(details.NumberOfSeasons), 10)

	if details.PosterPath != "" {
		tvshow.PosterURL = tmdb.GetImageURL(details.PosterPath, tmdb.Original)
	}

	if details.BackdropPath != "" {
		tvshow.BackdropURL = tmdb.GetImageURL(details.BackdropPath, tmdb.Original)
	}

	if details.TVImagesAppend != nil && details.TVImagesAppend.Images != nil {
		if len(details.Images.Logos) > 0 {
			tvshow.LogoURL = tmdb.GetImageURL(details.Images.Logos[0].FilePath, tmdb.Original)
		}
	}

	return
}

func (t *DatabaseTMDB) AskTvshowEpisode(q *EpisodeQuestion) (episode *TvshowEpisodeAnswer, err error) {
	episode = new(TvshowEpisodeAnswer)

	var (
		show_details *tmdb.TVDetails
		details      *tmdb.TVEpisodeDetails
	)

	show_details, err = t.get_tv_details(q.ShowID)
	if err != nil {
		return nil, err
	}

	details, err = t.Client.GetTVEpisodeDetails(q.ShowID, q.Season, q.Episode, map[string]string{
		"append_to_response": "credits",
	})
	if err != nil {
		return nil, err
	}

	tmdb_id := strconv.FormatInt(int64(details.ID), 10)

	episode.Title = details.Name
	if len(details.AirDate) > 4 {
		episode.Year = details.AirDate[0:4]
	}
	episode.Originaltitle = show_details.OriginalName
	episode.Showtitle = show_details.Name
	episode.Plot = details.Overview
	episode.Premiered = details.AirDate
	episode.ID = tmdb_id
	episode.Season = strconv.FormatInt(int64(details.SeasonNumber), 10)
	episode.Episode = strconv.FormatInt(int64(details.EpisodeNumber), 10)
	episode.Uniqueids = []nfo.ID{
		{
			Type:    "tmdb",
			Default: "true",
			Text:    tmdb_id,
		},
	}

	if details.Credits != nil {
		for _, credit := range details.Credits.Crew {
			switch credit.Job {
			case "Director":
				episode.Directors = append(episode.Directors, credit.Name)
			}
		}
	}

	if details.StillPath != "" {
		episode.ThumbURL = tmdb.GetImageURL(details.StillPath, tmdb.Original)
		episode.BackdropURL = tmdb.GetImageURL(details.StillPath, tmdb.Original)

		episode.Thumbs = append(episode.Thumbs, nfo.Thumb{
			Preview: tmdb.GetImageURL(details.StillPath, tmdb.W780),
			Text:    episode.ThumbURL,
		})
	}

	return
}
