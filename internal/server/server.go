package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/logging"
)

// Run creates and configures and run a new HTTP server.
// a is the address the server will listen on.
// h is the handler that will be used by router the incoming requests.
func Run(cfg *config.Config, h *handlers.URLHandler, logger *logging.Logger) (*http.Server, error) {
	srv := &http.Server{
		Addr:              cfg.Host,
		Handler:           Router(h),
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting server at %s, TLS enabled: %t", srv.Addr, cfg.EnableHTTPS))
		var err error
		if cfg.EnableHTTPS {
			err = srv.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", err)
		}
	}()

	return srv, nil
}
