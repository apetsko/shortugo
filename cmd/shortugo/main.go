package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/server"
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

	handler := handlers.NewURLHandler(cfg.BaseURL, storage, logger, cfg.Secret)

	// Batch deletion
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go storages.StartBatchDeleteProcessor(ctx, storage, handler.ToDelete, logger)

	// Start server
	srv, err := server.Run(cfg, handler, logger)
	if err != nil {
		logger.Fatal("Server failed: " + err.Error())
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()
	logger.Info("Shutting down server...")

	timeout := 5 * time.Second
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Graceful shutdown failed: " + err.Error())
	}

	logger.Info("Server gracefully stopped")
}
