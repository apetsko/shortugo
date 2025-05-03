// Package config provides functionality for managing application configuration.
// It supports loading configuration values from environment variables and command-line flags,
// ensuring flexibility and ease of use in different deployment environments.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/apetsko/shortugo/internal/utils"
	"github.com/caarlos0/env/v11"
)

// Config holds the configuration values for the application.
// These values can be populated from environment variables or command-line flags.
type Config struct {
	Config string `env:"CONFIG" envDefault:""`

	// Host is the network address with port for the server to listen on.
	Host string `env:"SERVER_ADDRESS" validate:"required"`

	// BaseURL is the base URL of the application, typically used for routing or API endpoints.
	BaseURL string `env:"BASE_URL" validate:"required"`

	// FileStoragePath is the path where the file storage is located.
	FileStoragePath string `env:"FILE_STORAGE_PATH" validate:"required_without=DatabaseDSN"`

	// DatabaseDSN is the Data Source Name (DSN) for connecting to the database.
	DatabaseDSN string `env:"DATABASE_DSN" validate:"required_without=FileStoragePath"`

	// Secret is the HMAC secret key used for signing and verifying data.
	Secret string `env:"SECRET" validate:"required"`

	// Cert is the file path to the SSL/TLS certificate used for HTTPS.
	TLSCertPath string `env:"CERT_FILE" validate:"required_if=EnableHTTPS true"`

	// Key is the file path to the SSL/TLS private key used for HTTPS.
	TLSKeyPath string `env:"KEY_FILE" validate:"required_if=EnableHTTPS true"`

	// Https indicates whether the application should use HTTPS for secure communication.
	EnableHTTPS bool `env:"ENABLE_HTTPS"`
}

// New creates a new Config instance, populating it with values from command-line flags and environment variables.
// Returns a pointer to the Config instance or an error if the configuration is invalid.
func New() (*Config, error) {
	var c Config

	// Parse command-line flags
	flag.BoolVar(&c.EnableHTTPS, "s", false, "enable https")
	flag.StringVar(&c.TLSCertPath, "cert", "certs/cert.crt", "certificate filepath")
	flag.StringVar(&c.TLSKeyPath, "key", "certs/cert.key", "private key filepath")
	flag.StringVar(&c.Config, "config", "", "private key filepath")
	flag.StringVar(&c.Host, "a", "localhost:8080", "network address with port")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "base url address")
	flag.StringVar(&c.FileStoragePath, "f", "db.json", "file storages name")
	flag.StringVar(&c.DatabaseDSN, "d", "", "database DSN")
	flag.StringVar(&c.Secret, "secret", "fortytwo", "HMAC secret")

	// Parse config.json
	if c.Config != "" {
		if err := LoadJSONConfig(c.Config, &c); err != nil {
			return nil, fmt.Errorf("failed to load config from json file: %w", err)
		}
	}

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

// LoadJSONConfig reads config.json file
func LoadJSONConfig(path string, out interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close config file: %s", err)
		}
	}()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("decode config: %w", err)
	}

	return nil
}
