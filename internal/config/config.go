package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Сonfig struct {
	Host    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func Parse() (c Сonfig, err error) {
	if err = env.Parse(c); err != nil {
		return
	}

	if c.BaseURL != "" && c.Host != "" {
		return
	}

	flag.StringVar(&c.Host, "a", "localhost:8080", "network address with port")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url address")
	flag.Parse()

	return c, nil
}
