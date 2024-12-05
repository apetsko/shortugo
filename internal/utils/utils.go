package utils

import (
	"errors"
	"fmt"
	"log"
)

func Generate(URL string) (ID string, err error) {
	ID = "EwHXdJfB"
	return
}

func FullURL(baseURL string, ID string) (string, error) {
	if ID == "" {
		return "", errors.New("empty id")
	}
	return fmt.Sprintf("%s/%s", baseURL, ID), nil
}
