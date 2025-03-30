package server

import (
	"github.com/apetsko/shortugo/internal/handlers"
	mw "github.com/apetsko/shortugo/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(handler *handlers.URLHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(mw.LogMiddleware(handler.Logger))
	r.Use(mw.GzipMiddleware(handler.Logger))

	r.Post("/", handler.ShortenURL)
	r.Post("/api/shorten", handler.ShortenJSON)
	r.Post("/api/shorten/batch", handler.ShortenBatchJSON)
	r.Get("/api/user/urls", handler.ListUserURLs)
	r.Delete("/api/user/urls", handler.DeleteUserURLs)
	r.Get("/{id}", handler.ExpandURL)
	r.Get("/ping", handler.PingDB)

	return r
}
