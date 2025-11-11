package server

import (
	"log"
	"os"
	"path/filepath"
)

func Main() {
	wd, _ := os.Getwd()

	configLocation := filepath.Join(wd, "server.conf")

	var s Server
	if err := s.config.Load(configLocation); err != nil {
		log.Fatal(err)
		return
	}

	if err := s.init(); err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(s.Serve())
}
