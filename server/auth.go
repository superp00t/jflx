package server

import (
	"net/http"
)

type auth_provider interface {
	AuthenticateCredentials(address, usergroup, username, password string) bool
}

func parse_address(s *Server, r *http.Request) string {
	// TODO
	return "127.0.0.1"
}

func (s *Server) authorize_request(group string, rw http.ResponseWriter, r *http.Request) (authorized bool) {
	if group == "" {
		authorized = true
		return
	}

	if s.config.AuthProvider == "" {
		return true
	}

	username, password, basic_auth := r.BasicAuth()
	if !basic_auth {
		rw.Header().Set("WWW-Authenticate", `Basic realm="jflx", charset="utf-8"`)
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !s.auth.AuthenticateCredentials(parse_address(s, r), group, username, password) {
		rw.Header().Set("WWW-Authenticate", `Basic realm="jflx", charset="utf-8"`)
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized = true
	return
}
