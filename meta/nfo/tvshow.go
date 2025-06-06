package nfo

import (
	"encoding/xml"
	"os"
)

type NamedSeason struct {
	Text   string `xml:",chardata"`
	Number string `xml:"number,attr"`
}

type Tvshow struct {
	XMLName        xml.Name      `xml:"tvshow"`
	Text           string        `xml:",chardata"`
	Title          string        `xml:"title"`
	Originaltitle  string        `xml:"originaltitle"`
	Showtitle      string        `xml:"showtitle"`
	Userrating     string        `xml:"userrating"`
	Top250         string        `xml:"top250"`
	Season         string        `xml:"season"`
	Episode        string        `xml:"episode"`
	Displayseason  string        `xml:"displayseason"`
	Displayepisode string        `xml:"displayepisode"`
	Outline        string        `xml:"outline"`
	Plot           string        `xml:"plot"`
	Tagline        string        `xml:"tagline"`
	Runtime        string        `xml:"runtime"`
	Mpaa           string        `xml:"mpaa"`
	Uniqueids      []ID          `xml:"uniqueid"`
	Genre          string        `xml:"genre"`
	Premiered      string        `xml:"premiered"`
	Year           string        `xml:"year"`
	Status         string        `xml:"status"`
	Code           string        `xml:"code"`
	Aired          string        `xml:"aired"`
	Studio         string        `xml:"studio"`
	Trailer        string        `xml:"trailer"`
	Actors         []Actor       `xml:"actor"`
	NamedSeason    []NamedSeason `xml:"namedseason"`
}

func WriteTvshow(filename string, m *Tvshow) error {
	b, err := xml.MarshalIndent(m, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0700)
}

func ReadTvshow(filename string, m *Tvshow) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, m)
}
