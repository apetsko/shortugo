package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/logging"
	grpch "github.com/apetsko/shortugo/internal/server/grpc/handlers"
	"github.com/apetsko/shortugo/internal/server/http/handlers"
	pb "github.com/apetsko/shortugo/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Run starts a gRPC server on the address specified in cfg.GRPCHost.
// If cfg.EnableHTTPS is true and TLSCertPath/TLSKeyPath are provided,
// the server will use TLS credentials.
//
// The function performs graceful shutdown on context cancellation,
// and logs startup and shutdown events asynchronously.
//
// Parameters:
//   - cfg: server configuration (host, TLS options)
//   - h: the URLHandler containing storage, logger, and auth logic
//   - logger: logger for lifecycle and error reporting
//
// Returns:
//   - *grpc.Server: the running gRPC server instance
//   - error: non-nil if the server fails to start
func Run(cfg *config.Config, h *handlers.URLHandler, logger *logging.Logger) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", cfg.GRPCHost)
	if err != nil {
		return nil, fmt.Errorf("listen error: %w", err)
	}

	var opts []grpc.ServerOption

	if cfg.EnableHTTPS {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertPath, cfg.TLSKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	srv := grpc.NewServer(opts...)
	pb.RegisterURLShortenerServer(srv, grpch.NewHandler(h))

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		<-ctx.Done()
		srv.GracefulStop()
		return nil
	})

	g.Go(func() error {
		logger.Info(fmt.Sprintf("Starting gRPC server at %s, TLS: %t", cfg.GRPCHost, cfg.EnableHTTPS))
		return srv.Serve(lis)
	})

	go func() {
		if err := g.Wait(); err != nil {
			logger.Error("gRPC server error", err)
		}
	}()

	return srv, nil
}
