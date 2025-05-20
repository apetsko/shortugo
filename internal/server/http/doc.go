// Package http provides the HTTP routing and middleware setup for the Shortugo URL shortening service.
// It defines the Router function that configures all RESTful endpoints, user authentication,
// gzip compression, request logging, panic recovery, and exposes internal diagnostic routes via pprof.
//
// The package integrates handlers for creating, expanding, deleting, and listing shortened URLs,
// as well as health checks and internal statistics. It also exposes profiling endpoints under /debug/pprof.
package http
