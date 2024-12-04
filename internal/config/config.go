package config

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type config struct {
	Host    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func Parse() config {
	c := new(config)

	if err := env.Parse(c); err == nil {
		if c.BaseURL != "" && c.Host != "" {
			return *c
		}

	} else {
		log.Fatal(err)
	}

	flag.StringVar(&c.Host, "a", "localhost:8080", "network address with port")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url address")
	flag.Parse()

	if c.BaseURL == "" || c.Host == "" {
		fmt.Fprintf(flag.CommandLine.Output(), "Wrong params: baseUrl: %q, host: %q", c.BaseURL, c.Host)
		flag.PrintDefaults()
		log.Fatalf("invalid parameters: Host and BaseURL are required")
	}
	return *c
}
