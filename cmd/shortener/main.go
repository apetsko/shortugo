package main

import (
	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/server"
	"github.com/apetsko/shortugo/internal/storage/inmem"
)

func main() {
	cfg := config.Parse()
	storage := inmem.New()
	handler := handlers.NewURLHandler(cfg.BaseURL, storage)

	router := handlers.SetupRouter(handler)

	server.StartServer(cfg.Host, router)
}
