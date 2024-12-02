package server

import (
	"net/http"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/storage/inmem"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Serve() {
	storage := inmem.New()
	handler := handlers.NewURLHandler(storage)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", handler.ShortenURL)
	r.Get("/{id}", handler.ExpandURL)

	http.ListenAndServe(config.Host, r)
}
