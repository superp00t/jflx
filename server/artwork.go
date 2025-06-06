package server

import (
	"os"
	"time"
)

type art_type struct {
	Width  int
	Height int
}

var poster_artwork = art_type{}

var default_artwork = art_type{}

var thumb_artwork = art_type{
	Width: 960,
}

var backdrop_artwork = art_type{}

var logo_artwork = art_type{}

func (art art_type) constrained() bool {
	return art.Width != 0 || art.Height != 0
}

func (s *Server) update_jpeg_file(art art_type, image_url string, image_file string) (err error) {
	fi, err := os.Stat(image_file)
	if err == nil {
		if time.Since(fi.ModTime()) < (24 * time.Hour) {
			return
		}
	}

	img, err := s.download_jpeg_url(art, image_url)
	if err == nil {
		err = os.WriteFile(
			image_file,
			img,
			0700)
	}

	return
}
