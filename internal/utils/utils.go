// Package utils provides utility functions and helpers used across the application.
// It includes functions for generating IDs, validating structs, and other reusable utilities.
package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validateInstance *validator.Validate
	once             sync.Once
)

// GenerateID generates a unique ID based on the SHA-256 hash of the input string.
// s is the input string to hash.
// length is the desired length of the generated ID.
func GenerateID(s string, length int) (id string) {
	hash := sha256.Sum256([]byte(s))
	id = base64.RawURLEncoding.EncodeToString(hash[:length])[:length]
	return
}

// GenerateUserID generates a random user ID of the specified length.
// length is the desired length of the generated user ID.
func GenerateUserID(length int) (id string, err error) {
	r := make([]byte, length)

	_, err = rand.Read(r)
	if err != nil {
		err = fmt.Errorf("failed to generate random User ID: %w", err)
		return "", err
	}

	id = hex.EncodeToString(r)

	return id, nil
}

func getValidator() *validator.Validate {
	once.Do(func() {
		validateInstance = validator.New()
	})
	return validateInstance
}

// ValidateStruct validates the fields of a struct based on the tags defined in the struct.
// a is the struct to validate.
func ValidateStruct(a any) error {
	return getValidator().Struct(a)
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
