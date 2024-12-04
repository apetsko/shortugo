package utils

import (
	"errors"
	"fmt"
	"log"
)

var baseURL string

func SetBaseUrl(u string) {
	baseURL = u
	log.Println("Base URL set successfully")
}

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
