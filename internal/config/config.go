package config

import (
	"flag"
	"fmt"
	"log"
)

type config struct {
	Host    string
	BaseURL string
}

func Parse() config {
	c := new(config)

	flag.StringVar(&c.Host, "Host", "localhost:8080", "network address with port")
	flag.StringVar(&c.BaseURL, "BaseURL", "localhost:8080", "base url address")

	flag.Parse()
	if c.BaseURL == "" || c.Host == "" {
		fmt.Fprintf(flag.CommandLine.Output(), "Wrong params: baseUrl: %q, host: %q", c.BaseURL, c.Host)
		flag.PrintDefaults()
		log.Fatalf("invalid parameters: Host and BaseURL are required")
	}
	log.Println("cccccccccccccc", c)
	return *c
}



