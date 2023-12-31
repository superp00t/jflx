package nfo

import (
	"encoding/xml"
	"io/ioutil"
)

type Movie struct {
	XMLName       xml.Name `xml:"movie"`
	Text          string   `xml:",chardata"`
	Title         string   `xml:"title"`
	Originaltitle string   `xml:"originaltitle"`
	Userrating    string   `xml:"userrating"`
	Plot          string   `xml:"plot"`
	Mpaa          string   `xml:"mpaa"`
	Uniqueid      *ID      `xml:"uniqueid"`
	Genre         string   `xml:"genre"`
	Tagline       string   `xml:"tagline"`
	Tag           string   `xml:"tag"`
	Country       string   `xml:"country"`
	Credits       string   `xml:"credits"`
	Director      string   `xml:"director"`
	Premiered     string   `xml:"premiered"`
	Studio        string   `xml:"studio"`
	Actors        []Actor  `xml:"actor"`
}

func (m *Movie) Write(filename string) error {
	b, err := xml.MarshalIndent(m, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0700)
}

func (m *Movie) Read(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, m)
}
