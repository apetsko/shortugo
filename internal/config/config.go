package config

import (
	"flag"
	"fmt"

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
	flag.StringVar(&c.FileStoragePath, "f", "db.json", "file storages name")
	flag.StringVar(&c.DatabaseDSN, "d", "", "database DSN")

	flag.Parse()

	if err = env.Parse(&c); err != nil {
		return c, fmt.Errorf("error while parse envs: %w", err)
	}

	return c, nil
}
