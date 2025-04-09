package config

import (
	"flag"
	"fmt"

	"github.com/apetsko/shortugo/internal/utils"
	"github.com/caarlos0/env/v11"
)

// Config holds the configuration values for the application.
// These values can be populated from environment variables or command-line flags.
type Config struct {
	// Host is the network address with port for the server to listen on.
	Host string `env:"SERVER_ADDRESS"`

	// BaseURL is the base URL of the application, typically used for routing or API endpoints.
	BaseURL string `env:"BASE_URL"`

	// FileStoragePath is the path where the file storage is located.
	FileStoragePath string `env:"FILE_STORAGE_PATH"`

	// DatabaseDSN is the Data Source Name (DSN) for connecting to the database.
	DatabaseDSN string `env:"DATABASE_DSN"`

	// Secret is the HMAC secret key used for signing and verifying data.
	Secret string `env:"SECRET"`
}

// New creates a new Config instance, populating it with values from command-line flags and environment variables.
// Returns a pointer to the Config instance or an error if the configuration is invalid.
func New() (*Config, error) {
	var c Config

	// Parse command-line flags
	flag.StringVar(&c.Host, "a", "localhost:8080", "network address with port")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url address")
	flag.StringVar(&c.FileStoragePath, "f", "db.json", "file storages name")
	flag.StringVar(&c.DatabaseDSN, "d", "", "database DSN")
	flag.StringVar(&c.Secret, "s", "fortytwo", "HMAC secret")

	// Parse the flags
	flag.Parse()

	// Load environment variables into the Config struct
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to load environment: %w", err)
	}

	// Validate the loaded configuration
	if err := utils.ValidateStruct(c); err != nil {
		return nil, err
	}

	// Return the populated and validated Config
	return &c, nil
}
