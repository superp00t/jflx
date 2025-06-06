package meta

import "github.com/superp00t/jflx/meta/nfo"

type ShowQuestion struct {
	Name   string
	Year   int
	HintID int
}

type EpisodeQuestion struct {
	ShowID           int
	Season           int
	Episode          int
	HintEpisodeTitle string
}

type Artwork struct {
	PosterURL   string
	BackdropURL string
	ThumbURL    string
	LogoURL     string
}

type MovieAnswer struct {
	nfo.Movie
	Artwork
}

type TvshowAnswer struct {
	nfo.Tvshow
	Artwork
}

type TvshowEpisodeAnswer struct {
	nfo.TvshowEpisode
	Artwork
}

type Source interface {
	AskMovie(q *ShowQuestion) (*MovieAnswer, error)
	AskTvshow(q *ShowQuestion) (*TvshowAnswer, error)
	AskTvshowEpisode(q *EpisodeQuestion) (*TvshowEpisodeAnswer, error)
}
