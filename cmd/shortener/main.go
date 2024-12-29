package main

import (
	"log"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/server"
	"github.com/apetsko/shortugo/internal/storage/inmem"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}

	storage := inmem.New()

	handler := handlers.NewURLHandler(cfg.BaseURL, storage)
	router := handlers.SetupRouter(handler)

	server := server.New(cfg.BaseURL, router)
	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}
}
