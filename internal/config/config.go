package config

import (
	"flag"
	"fmt"

	"github.com/apetsko/shortugo/internal/utils"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host            string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Secret          string `env:"SECRET"`
}

func New() (*Config, error) {
	var c Config

	flag.StringVar(&c.Host, "a", "localhost:8080", "network address with port")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url address")
	flag.StringVar(&c.FileStoragePath, "f", "db.json", "file storages name")
	flag.StringVar(&c.DatabaseDSN, "d", "", "database DSN")
	flag.StringVar(&c.Secret, "s", "fortytwo", "HMAC secret")

	flag.Parse()

	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to load environment: %w", err)
	}

	if err := utils.ValidateStruct(c); err != nil {
		return nil, err
	}
	return &c, nil
}
