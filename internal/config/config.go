package config

import (
	"flag"
	"fmt"

	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/storage/infile"
	"github.com/apetsko/shortugo/internal/storage/inmem"
	"github.com/apetsko/shortugo/internal/storage/postgres"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host            string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func Parse() (c Config, err error) {
	flag.StringVar(&c.Host, "a", "localhost:8080", "network address with port")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url address")
	flag.StringVar(&c.FileStoragePath, "f", "db.json", "file storage name")
	flag.StringVar(&c.DatabaseDSN, "d", "", "database DSN")

	flag.Parse()

	if err = env.Parse(&c); err != nil {
		return c, fmt.Errorf("error while parse envs: %w", err)
	}

	return c, nil
}

func (c Config) Storage(logger *logging.ZapLogger) (s handlers.Storage, err error) {
	switch {
	case c.DatabaseDSN != "":
		s, err = postgres.New(c.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		logger.Info("Using database storage")
		return s, nil
	case c.FileStoragePath != "":
		s, err = infile.New(c.FileStoragePath)
		if err != nil {
			return nil, err
		}
		logger.Info("Using file storage")
		return s, nil
	default:
		s = inmem.New()
		logger.Info("Using in-memory storage")
		return s, nil
	}
}
