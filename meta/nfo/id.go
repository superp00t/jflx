package nfo

type ID struct {
	Text    string `xml:",chardata"`
	Type    string `xml:"type,attr"`
	Default string `xml:"default,attr"`
}

func GetDefault(ids []ID) (id *ID) {
	for i := range ids {
		id = &ids[i]
		if id.Default == "true" {
			return
		}
	}

	id = nil
	return
}
