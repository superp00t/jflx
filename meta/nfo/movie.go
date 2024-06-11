package nfo

import (
	"encoding/xml"
	"os"
)

type Movie struct {
	XMLName       xml.Name `xml:"movie"`
	Text          string   `xml:",chardata"`
	Title         string   `xml:"title"`
	Originaltitle string   `xml:"originaltitle"`
	Userrating    string   `xml:"userrating"`
	Plot          string   `xml:"plot"`
	Mpaa          string   `xml:"mpaa"`
	Uniqueids     []ID     `xml:"uniqueid"`
	Genre         string   `xml:"genre"`
	Tagline       string   `xml:"tagline"`
	Tag           string   `xml:"tag"`
	Country       string   `xml:"country"`
	Credits       string   `xml:"credits"`
	Directors     []string `xml:"director"`
	Premiered     string   `xml:"premiered"`
	Studio        string   `xml:"studio"`
	Actors        []Actor  `xml:"actor"`
}

func WriteMovie(filename string, m *Movie) error {
	b, err := xml.MarshalIndent(m, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0700)
}

func ReadMovie(filename string, m *Movie) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, m)
}
