package server

import (
	"log"
	"os"
	"path/filepath"

	"github.com/superp00t/jflx/conf"
)

func exeDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func Main() {
	configLocation := filepath.Join(exeDir(), "server.conf")

	cfg, err := conf.LoadServer(configLocation)
	if err != nil {
		log.Fatal(err)
		return
	}

	s := new(Server)
	if err := s.Init(cfg); err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(s.Serve())
}
