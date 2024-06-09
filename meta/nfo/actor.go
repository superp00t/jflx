package nfo

type Actor struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name"`
	Role  string `xml:"role"`
	Order string `xml:"order"`
	Thumb string `xml:"thumb"`
}
