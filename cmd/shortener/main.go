package main

import (
	"log"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/server"
	"github.com/apetsko/shortugo/internal/storage/inmem"
)

func main() {
	logger, err := logging.NewZapLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	cfg, err := config.Parse()
	if err != nil {
		logger.Fatal(err.Error())
	}

	storage := inmem.New()

	handler := handlers.NewURLHandler(cfg.BaseURL, storage, logger)
	router := handlers.SetupRouter(handler)

	s := server.New(cfg.Host, router)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatal(err.Error())
	}
}
