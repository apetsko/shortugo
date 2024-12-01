package main

import (
	"net/http"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", handlers.ShortenURL)
	mux.HandleFunc("GET /{id}", handlers.ExpandURL)

	http.ListenAndServe(config.Host, mux)
}
