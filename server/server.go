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
	"github.com/superp00t/jflx/media/httpdirectory"
	"github.com/superp00t/jflx/meta"
)

type Server struct {
	// the configuration
	config conf.Server
	// http path router
	serve_mux *http.ServeMux
	server    *http.Server
	scraper   meta.Source
	// signals whether scraper is active
	scraper_status atomic.Bool
	// volumes
	volumes map[string]*Volume
	// auth
	auth auth_provider
}

func (s *Server) init_volumes() {
	s.volumes = make(map[string]*Volume)

	for i := range s.config.Volumes {
		volume := new(Volume)
		volume.config = s.config.Volumes[i]
		if volume.config.Cache != "" {
			cache_server, err := cache.NewServer(&cache.Config{
				Directory:       volume.config.Cache,
				MaxAge:          24 * time.Hour,
				MaxDirectoryAge: 2 * time.Minute,
				MaxSize:         volume.config.MaxCacheSize,
			}, volume)
			if err != nil {
				log.Fatal(err)
			}
			volume.handler = httpdirectory.FileServer(cache_server)
		} else {
			volume.handler = httpdirectory.FileServer(volume)
		}

		s.volumes[volume.config.Handle] = volume
	}
}

func (s *Server) handle_volume(rw http.ResponseWriter, r *http.Request) {
	volume_handle := r.PathValue("volume")

	volume, found := s.volumes[volume_handle]
	if !found {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	if volume.config.UserGroup != "" {
		if !s.authorize_request(volume.config.UserGroup, rw, r) {
			return
		}
	}

	http.StripPrefix("/media/"+volume_handle, volume.handler).ServeHTTP(rw, r)
}

func (s *Server) handle_get_list_volumes(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == "HEAD" {
		return
	}

	var ks []string
	for _, v := range s.config.Volumes {
		if !v.Unlisted {
			ks = append(ks, v.Handle)
		}
	}
	sort.Strings(ks)

	fmt.Fprintf(rw, "<h1 style=\"font-family: Arial;\">All JFLX Volumes</h1><pre>\n")
	for _, k := range ks {
		volume := s.volumes[k]
		if !volume.config.Unlisted {
			fmt.Fprintf(rw, "<a href=\"%s/\">%s (%s)</a>\n", volume.config.Handle, volume.config.Handle, volume.config.Kinds.String())
		}
	}
	fmt.Fprintf(rw, "</pre>")
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println(r.Header.Get("Referer"), fmt.Sprintf("HTTP/%d.%d", r.ProtoMajor, r.ProtoMinor), r.RemoteAddr, r.Method, r.URL.Path)
	s.serve_mux.ServeHTTP(rw, r)
}

func (s *Server) init() (err error) {
	switch s.config.AuthProvider {
	case "ldap":
		s.auth = new_ldap_cached_auth_provider(s)
	case "":
	default:
		err = fmt.Errorf("unknown auth provider: %s", s.config.AuthProvider)
		return
	}

	s.init_volumes()

	s.serve_mux = http.NewServeMux()
	s.serve_mux.HandleFunc("POST /api/v1/refresh", s.handle_post_refresh)
	s.serve_mux.HandleFunc("/media/", s.handle_get_list_volumes)
	s.serve_mux.HandleFunc("/media/{volume}/", s.handle_volume)

	s.server = &http.Server{
		Handler:      s,
		Addr:         s.config.ListenAddress,
		WriteTimeout: 5 * time.Hour,
		ReadTimeout:  5 * time.Hour,
	}

	return
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}
