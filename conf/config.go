package conf

import (
	"fmt"
	"strconv"
	"strings"
)

type MediaKinds uint16

const (
	Movie = 1 << iota
	TvShow
	Music
	Pictures
	Books
	Podcast
)

var kndtable = map[string]MediaKinds{
	"file":  0,
	"movie": Movie,
	"tv":    TvShow,
	"music": Music,
	"pics":  Pictures,
	"books": Books,
	"pod":   Podcast,
}

var kindtable = map[MediaKinds]string{
	Movie:    "Films",
	TvShow:   "TV Programs",
	Music:    "Music",
	Pictures: "Images",
	Books:    "Books",
	Podcast:  "Podcasts",
}

func (mks MediaKinds) String() string {
	var ks []string

	for b := 0; b < 16; b++ {
		mask := MediaKinds(1 << b)
		if mks&mask != 0 {
			ks = append(ks, kindtable[mask])
		}
	}

	return strings.Join(ks, ", ")
}

func (mks *MediaKinds) UnmarshalJSON(b []byte) error {
	str, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	exploded := strings.Split(str, ",")

	for _, kind := range exploded {
		bit, ok := kndtable[kind]
		if ok {
			*mks |= bit
		} else {
			return fmt.Errorf("conf: (*MediaKinds).UnmarshalJSON: Unknown media kind %s", kind)
		}
	}

	return nil
}

type Volume struct {
	Kinds MediaKinds
	// Handle: determines what the URL prefix will be i.e. (/media/{Handle}/movie.mp4)
	Handle string
	// Sources: the
	Sources []string
	// Cache: The directory to store a cache.
	Cache string
	// MaxCacheSize: the amount of bytes the directory is not allowed to exceed
	MaxCacheSize int64

	Unlisted bool
}

type Server struct {
	ListenAddress string
	TMDBScrapeKey string
	Volumes       []Volume
	Tokens        []string
}
