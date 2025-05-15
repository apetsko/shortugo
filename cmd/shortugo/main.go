package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/server/grpc"
	"github.com/apetsko/shortugo/internal/server/http"
	"github.com/apetsko/shortugo/internal/server/http/handlers"
	"github.com/apetsko/shortugo/internal/storages"
	"go.uber.org/zap/zapcore"
)

// Build info vars
var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	fmt.Println("Build version: " + BuildVersion)
	fmt.Println("Build date: " + BuildDate)
	fmt.Println("Build commit: " + BuildCommit)

	logger, err := logging.New(zapcore.DebugLevel)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	logger.Infof("Starting server with LogLevel: %s", zapcore.DebugLevel)

	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err.Error())
	}

	storage, err := storages.Init(cfg.DatabaseDSN, cfg.FileStoragePath, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer func() {
		err = storage.Close()
		if err != nil {
			logger.Fatal("failed to close storage: " + err.Error())
		}
	}()

	handler := handlers.NewURLHandler(cfg.BaseURL, storage, logger, cfg.Secret, cfg.TrustedSubnet)

	// Batch deletion
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go storages.StartBatchDeleteProcessor(ctx, storage, handler.ToDelete, logger)

	// Start HTTP server
	if _, err := http.Run(cfg, handler, logger); err != nil {
		logger.Fatal("HTTP server failed: " + err.Error())
	}

	// Start gRPC server
	if _, err := grpc.Run(cfg, handler, logger); err != nil {
		logger.Fatal("gRPC server failed: " + err.Error())
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()
	logger.Info("Shutting down servers...")
}
