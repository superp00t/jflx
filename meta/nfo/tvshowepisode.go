package nfo

import (
	"encoding/xml"
	"os"
)

type Thumb struct {
	Text    string `xml:",chardata"`
	Spoof   string `xml:"spoof,attr"`
	Cache   string `xml:"cache,attr"`
	Aspect  string `xml:"aspect,attr"`
	Preview string `xml:"preview,attr"`
}

type TvshowEpisode struct {
	XMLName        xml.Name `xml:"episodedetails"`
	Text           string   `xml:",chardata"`
	Title          string   `xml:"title"`
	Originaltitle  string   `xml:"originaltitle"`
	Showtitle      string   `xml:"showtitle"`
	Userrating     string   `xml:"userrating"`
	Top250         string   `xml:"top250"`
	Season         string   `xml:"season"`
	Episode        string   `xml:"episode"`
	Displayseason  string   `xml:"displayseason"`
	Displayepisode string   `xml:"displayepisode"`
	Outline        string   `xml:"outline"`
	Plot           string   `xml:"plot"`
	Tagline        string   `xml:"tagline"`
	Runtime        string   `xml:"runtime"`
	Thumbs         []Thumb  `xml:"thumb"`
	Mpaa           string   `xml:"mpaa"`
	Playcount      string   `xml:"playcount"`
	Lastplayed     string   `xml:"lastplayed"`
	ID             string   `xml:"id"`
	Uniqueids      []ID     `xml:"uniqueid"`
	Genre          []string `xml:"genre"`
	Credits        string   `xml:"credits"`
	Directors      []string `xml:"director"`
	Premiered      string   `xml:"premiered"`
	Year           string   `xml:"year"`
	Status         string   `xml:"status"`
	Code           string   `xml:"code"`
	Aired          string   `xml:"aired"`
	Studio         string   `xml:"studio"`
	Trailer        string   `xml:"trailer"`
	Actors         []Actor  `xml:"actor"`
	Dateadded      string   `xml:"dateadded"`
}

func WriteTvshowEpisode(filename string, m *TvshowEpisode) error {
	b, err := xml.MarshalIndent(m, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0700)
}

func ReadTvshowEpisode(filename string, m *TvshowEpisode) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, m)
}
