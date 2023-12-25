package nfo

type ID struct {
	Text    string `xml:",chardata"`
	Type    string `xml:"type,attr"`
	Default string `xml:"default,attr"`
}
