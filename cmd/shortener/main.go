package main

import (
	"net/http"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/storage/inmem"
)

func main() {
	storage := inmem.NewInMem()
	handler := handlers.NewURLHandler(storage)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", handler.ShortenURL)
	mux.HandleFunc("GET /{id}", handler.ExpandURL)

	http.ListenAndServe(config.Host, mux)
}
