package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/mocks"
	"github.com/apetsko/shortugo/internal/server/http/handlers"
	pb "github.com/apetsko/shortugo/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCServer_Starts(t *testing.T) {
	// Setup mock storage
	mockStorage := new(mocks.Storage)
	mockStorage.On("Ping").Return(nil)

	logger, _ := logging.New(zapcore.DebugLevel)

	addr := "127.0.0.1:19090"

	cfg := &config.Config{
		GRPCHost: addr,
	}

	urlHandler := handlers.NewURLHandler("http://localhost", mockStorage, logger, "secret", "127.0.0.0/8")

	srv, err := Run(cfg, urlHandler, logger)
	require.NoError(t, err)
	require.NotNil(t, srv)

	conn, err := grpc.NewClient(
		"passthrough:///"+addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer func() {
		err = conn.Close()
		require.NoError(t, err)
	}()

	client := pb.NewURLShortenerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Ping(ctx, &pb.PingRequest{})
	require.NoError(t, err)
	assert.Equal(t, "OK", resp.GetStatus())

	srv.GracefulStop()
}

func TestGRPCServer_TLSConfigError(t *testing.T) {
	mockStorage := new(mocks.Storage)
	mockStorage.On("Ping").Return(nil)

	logger, _ := logging.New(zapcore.DebugLevel)

	cfg := &config.Config{
		GRPCHost:    "127.0.0.1:19091",
		EnableHTTPS: true,
		TLSCertPath: "bad-cert.pem",
		TLSKeyPath:  "bad-key.pem",
	}

	urlHandler := handlers.NewURLHandler("http://localhost", mockStorage, logger, "secret", "127.0.0.0/8")

	srv, err := Run(cfg, urlHandler, logger)

	assert.Nil(t, srv)
	require.ErrorContains(t, err, "failed to load TLS credentials")
}
