package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	address string
	router  *chi.Mux
}

func New(a string, r *chi.Mux) *Server {
	return &Server{a, r}
}

func (s *Server) StartServer() error {
	log.Printf("Starting server on %s", s.address)
	return http.ListenAndServe(s.address, s.router)
}
