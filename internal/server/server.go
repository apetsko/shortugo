package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// New creates and configures a new HTTP server.
// a is the address the server will listen on.
// r is the router that will handle the incoming requests.
func New(a string, r *chi.Mux) *http.Server {
	server := &http.Server{
		Addr:              a,               // Address to listen on.
		ReadHeaderTimeout: 3 * time.Second, // Maximum duration for reading the request headers.
		Handler:           r,               // HTTP handler to invoke.
	}
	return server
}
