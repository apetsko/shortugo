package main

import (
	"log"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/server"
	"github.com/apetsko/shortugo/internal/storage/infile"
)

func main() {
	zlogger, err := logging.NewZapLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	cfg, err := config.Parse()
	if err != nil {
		zlogger.Fatal(err.Error())
	}

	storage, err := infile.New(cfg.FileStoragePath)
	if err != nil {
		zlogger.Fatal(err.Error())
	}

	handler := handlers.NewURLHandler(cfg.BaseURL, storage, zlogger)
	router := handlers.SetupRouter(handler)
	s := server.New(cfg.Host, router)

	zlogger.Info("running server on " + cfg.Host)
	if err := s.ListenAndServe(); err != nil {
		zlogger.Fatal(err.Error())
	}
}
