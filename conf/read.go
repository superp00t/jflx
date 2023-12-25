package conf

import (
	"encoding/json"
	"io/ioutil"
)

func LoadServer(path string) (*Server, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := new(Server)
	return s, json.Unmarshal(b, s)
}
