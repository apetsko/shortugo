package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/server/http/handlers"
	"golang.org/x/sync/errgroup"
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

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		<-ctx.Done()
		five := 5 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), five)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	})

	g.Go(func() error {
		logger.Info(fmt.Sprintf("Starting HTTP server at %s, TLS: %t", srv.Addr, cfg.EnableHTTPS))
		if cfg.EnableHTTPS {
			return srv.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath)
		}
		return srv.ListenAndServe()
	})

	go func() {
		if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", err)
		}
	}()

	return srv, nil
}
