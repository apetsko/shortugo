package main

import (
	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/server"
	"github.com/apetsko/shortugo/internal/storage/inmem"
	"github.com/apetsko/shortugo/internal/utils"
)

func main() {
	cfg := config.Parse()

	utils.SetBaseURL(cfg.BaseURL)

	storage := inmem.New()
	handler := handlers.NewURLHandler(cfg.BaseURL, storage)

	router := handlers.SetupRouter(handler)

	server.StartServer(cfg.Host, router)
}
