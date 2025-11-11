package config

import (
	"encoding/json"
	"os"
)

func (server_config *ServerConfig) Load(name string) (err error) {
	var b []byte
	b, err = os.ReadFile(name)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, server_config)
	return
}
