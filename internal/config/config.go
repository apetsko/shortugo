package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Сonfig struct {
	Host            string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func Parse() (c Сonfig, err error) {
	flag.StringVar(&c.Host, "a", "localhost:8080", "network address with port")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url address")
	flag.StringVar(&c.FileStoragePath, "f", "db.json", "file storage name")
	flag.Parse()

	if err = env.Parse(&c); err != nil {
		return
	}

	return c, nil
}
