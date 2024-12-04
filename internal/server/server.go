package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StartServer(address string, router *chi.Mux) {
	log.Printf("Starting server on %s", address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("server error: %s", err.Error())
	}
}
