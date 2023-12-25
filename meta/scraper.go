package meta

import "github.com/superp00t/jflx/meta/nfo"

type MovieQuestion struct {
	Name   string
	Year   int
	HintID int
}

type MovieAnswer struct {
	nfo.Movie
	PosterURL   string
	BackdropURL string
	LogoURL     string
}

type Source interface {
	AskMovie(q *MovieQuestion) (*MovieAnswer, error)
}
