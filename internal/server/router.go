package server

import (
	"expvar"
	"net/http/pprof"

	"github.com/apetsko/shortugo/internal/handlers"
	mw "github.com/apetsko/shortugo/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router initializes the router with all the necessary routes and middleware.
func Router(handler *handlers.URLHandler) *chi.Mux {
	r := chi.NewRouter()
	// Middleware to get the real IP address of the client.
	r.Use(middleware.RealIP)
	// Middleware to recover from panics and return a 500 error.
	r.Use(middleware.Recoverer)
	// Custom middleware to log the details of each request and response.
	r.Use(mw.LogMiddleware(handler.Logger))
	// Custom middleware to compress the response body using gzip.
	r.Use(mw.GzipMiddleware(handler.Logger))

	// Route to shorten a URL.
	r.Post("/", handler.ShortenURL)
	// Route to shorten a URL via JSON request.
	r.Post("/api/shorten", handler.ShortenJSON)
	// Route to shorten multiple URLs via batch JSON request.
	r.Post("/api/shorten/batch", handler.ShortenBatchJSON)
	// Route to list all URLs associated with a user.
	r.Get("/api/user/urls", handler.ListUserURLs)
	// Route to delete multiple URLs associated with a user.
	r.Delete("/api/user/urls", handler.DeleteUserURLs)
	// Route to expand a shortened URL.
	r.Get("/{id}", handler.ExpandURL)
	// Route to check the database connection.
	r.Get("/ping", handler.PingDB)
	// Route to list all URLs associated with a user.
	r.Get("/api/internal/stats", handler.Stats)

	r.Route("/debug/pprof", func(r chi.Router) {
		r.HandleFunc("/*", pprof.Index)
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
		r.Handle("/vars", expvar.Handler())

		r.Handle("/goroutine", pprof.Handler("goroutine"))
		r.Handle("/threadcreate", pprof.Handler("threadcreate"))
		r.Handle("/mutex", pprof.Handler("mutex"))
		r.Handle("/heap", pprof.Handler("heap"))
		r.Handle("/block", pprof.Handler("block"))
		r.Handle("/allocs", pprof.Handler("allocs"))
	})

	return r
}
