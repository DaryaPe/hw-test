package internalhttp

import (
	"net/http"
)

func (s *Server) CreateMux() {
	mux := http.NewServeMux()
	mux.Handle("/health", s.loggingMiddleware(s.defaultRoute()))
	s.handler = mux
}
