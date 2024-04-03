package server

import (
	"fmt"
	"net/http"
)

func (s *Server) check_authorization(rw http.ResponseWriter, r *http.Request) (err error) {
	user_token := r.Header.Get("X-JFLX-Token")
	for _, token := range s.Conf.Tokens {
		if token == user_token {
			return
		}
	}
	err = fmt.Errorf("unauthorized")
	http.Error(rw, err.Error(), http.StatusUnauthorized)
	return
}

func (s *Server) handle_post_refresh(rw http.ResponseWriter, r *http.Request) {
	if err := s.check_authorization(rw, r); err != nil {
		return
	}

	if !s.scrape_in_progress() {
		fmt.Fprintf(rw, "starting")
		go s.perform_scrape_and_log()
	} else {
		fmt.Fprintf(rw, "scrape in progress")
	}
}
