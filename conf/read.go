package conf

import (
	"encoding/json"
	"os"
)

func LoadServer(name string, c *Server) (err error) {
	var b []byte
	b, err = os.ReadFile(name)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, c)
	return
}
