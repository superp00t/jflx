package server

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync/atomic"
	"time"

	"github.com/superp00t/jflx/cache"
	"github.com/superp00t/jflx/conf"
	"github.com/superp00t/jflx/media"
	"github.com/superp00t/jflx/meta"
)

type Server struct {
	Conf           *conf.Server
	Router         *http.ServeMux
	WebServer      *http.Server
	Volumes        map[string]*Volume
	scraper        meta.Source
	scraper_status atomic.Bool
}

func (s *Server) LoadVolumes() {
	s.Volumes = make(map[string]*Volume)

	for i := range s.Conf.Volumes {
		cvol := &s.Conf.Volumes[i]
		vol := new(Volume)
		vol.Conf = cvol
		if cvol.Cache != "" {
			cache_server, err := cache.NewServer(&cache.Config{
				Directory:       cvol.Cache,
				MaxAge:          24 * time.Hour,
				MaxDirectoryAge: 2 * time.Minute,
				MaxSize:         cvol.MaxCacheSize,
			}, vol)
			if err != nil {
				log.Fatal(err)
			}
			vol.Handler = media.FileServer(cache_server)
		} else {
			vol.Handler = media.FileServer(vol)
		}

		volumePrefix := fmt.Sprintf("/media/%s/", vol.Conf.Handle)

		s.Router.Handle(volumePrefix, http.StripPrefix(volumePrefix, vol.Handler))

		s.Volumes[vol.Conf.Handle] = vol
	}
}

func (s *Server) handle_get_list_volumes(rw http.ResponseWriter, r *http.Request) {
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
	s.Router = http.NewServeMux()
	s.Router.HandleFunc("POST /api/v1/refresh", s.handle_post_refresh)
	s.Router.HandleFunc("/media/", s.handle_get_list_volumes)

	s.LoadVolumes()

	s.WebServer = &http.Server{
		Handler:      s,
		Addr:         s.Conf.ListenAddress,
		WriteTimeout: 5 * time.Hour,
		ReadTimeout:  5 * time.Hour,
	}
	return nil
}

func (s *Server) Serve() error {
	return s.WebServer.ListenAndServe()
}
