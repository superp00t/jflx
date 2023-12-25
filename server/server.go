package server

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"github.com/superp00t/jflx/conf"
	"github.com/superp00t/jflx/media"
	"github.com/superp00t/jflx/meta"
)

type Server struct {
	Conf      *conf.Server
	Router    *mux.Router
	WebServer *http.Server
	Volumes   map[string]*Volume
	Scraper   meta.Source
}

func (s *Server) LoadVolumes() {
	s.Volumes = make(map[string]*Volume)

	for i := range s.Conf.Volumes {
		cvol := &s.Conf.Volumes[i]
		vol := new(Volume)
		vol.Conf = cvol
		vol.Handler = media.FileServer(vol)

		volumePrefix := fmt.Sprintf("/media/%s/", vol.Conf.Handle)

		route := s.Router.PathPrefix(volumePrefix)
		route.Handler(http.StripPrefix(volumePrefix, vol.Handler))

		s.Volumes[vol.Conf.Handle] = vol
	}
}

func (s *Server) ListVolumes(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == "HEAD" {
		return
	}

	var ks []string
	for k := range s.Volumes {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	fmt.Fprintf(rw, "<h1 style=\"font-family: Arial;\">All JFLX Volumes</h1><pre>\n")
	for _, k := range ks {
		vol := s.Volumes[k]
		if !vol.Conf.Unlisted {
			fmt.Fprintf(rw, "<a href=\"%s/\">%s (%s)</a>\n", vol.Conf.Handle, vol.Conf.Handle, vol.Conf.Kinds.String())
		}
	}
	fmt.Fprintf(rw, "</pre>")
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.URL.Path)
	s.Router.ServeHTTP(rw, r)
}

func (s *Server) Init(cfg *conf.Server) error {
	s.Conf = cfg
	r := mux.NewRouter()
	s.Router = r
	r.HandleFunc("/refresh", s.Refresh)
	r.HandleFunc("/media/", s.ListVolumes)

	s.LoadVolumes()

	s.WebServer = &http.Server{
		Handler: s,
		Addr:    s.Conf.ListenAddress,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return nil
}

func (s *Server) Serve() error {
	return s.WebServer.ListenAndServe()
}
