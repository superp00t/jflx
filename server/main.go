package server

import (
	"log"
	"os"
	"path/filepath"

	"github.com/superp00t/jflx/conf"
)

func Main() {
	wd, _ := os.Getwd()

	configLocation := filepath.Join(wd, "server.conf")

	s := new(Server)
	err := conf.LoadServer(configLocation, &s.config)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := s.init(); err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(s.Serve())
}
