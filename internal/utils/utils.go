package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

func Generate(URL string) (ID string, err error) {
	hash := sha256.Sum256([]byte(URL))
	ID = base64.RawURLEncoding.EncodeToString(hash[:6])
	return
}

func FullURL(baseURL string, ID string) (string, error) {
	if ID == "" {
		return "", errors.New("empty id")
	}
	return fmt.Sprintf("%s/%s", baseURL, ID), nil
}
