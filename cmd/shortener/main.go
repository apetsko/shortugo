package main

import (
	"log"

	"github.com/apetsko/shortugo/internal/config"
	"github.com/apetsko/shortugo/internal/handlers"

	zl "github.com/apetsko/shortugo/internal/log"
	"github.com/apetsko/shortugo/internal/server"
	"github.com/apetsko/shortugo/internal/storage/inmem"
)

func main() {
	err := zl.Start()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Parse()
	if err != nil {
		zl.Fatal(err.Error())
	}

	storage := inmem.New()

	handler := handlers.NewURLHandler(cfg.BaseURL, storage)
	router := handlers.SetupRouter(handler)

	s := server.New(cfg.Host, router)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
