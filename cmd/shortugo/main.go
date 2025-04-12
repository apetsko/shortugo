package main

import (
	"context"
	"fmt"
	"log"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/storages"
	"go.uber.org/zap/zapcore"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/server"
)

// Build info vars
var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	logger, err := logging.New(zapcore.DebugLevel)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	logger.Infof("Starting server with LogLevel: %s", zapcore.DebugLevel)

	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err.Error())
	}

	var storage handlers.Storage
	storage, err = storages.Init(cfg.DatabaseDSN, cfg.FileStoragePath, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer func(storage handlers.Storage) {
		err := storage.Close()
		if err != nil {
			logger.Fatal("failed to close storage: " + err.Error())
		}
	}(storage)

	handler := handlers.NewURLHandler(cfg.BaseURL, storage, logger, cfg.Secret)

	go storages.StartBatchDeleteProcessor(context.Background(), storage, handler.ToDelete, logger)

	router := server.Router(handler)
	s := server.New(cfg.Host, router)

	logger.Info("running server on " + cfg.Host)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatal(err.Error())
	}
}
