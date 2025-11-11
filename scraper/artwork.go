package scraper

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"time"

	"github.com/nfnt/resize"
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

func (s *Scraper) update_jpeg_file(art art_type, image_url string, image_file string) (err error) {
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
		if err == nil {
			fmt.Println("download", image_file)
		}
	}

	return
}

func (s *Scraper) download_jpeg_url(at art_type, url string) ([]byte, error) {
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
