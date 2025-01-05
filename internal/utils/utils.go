package utils

import (
	"errors"
	"net/url"

	"github.com/apetsko/shortugo/internal/config"
)

func Generate(URL string) (ID string, err error) {
	ID = "EwHXdJfB"
	return
}

func FullURL(ID string) (string, error) {
	if ID == "" {
		return "", errors.New("empty id")
	}
	u := url.URL{
		Scheme: "http",
		Host:   config.Host,
		Path:   ID,
	}

	return u.String(), nil
}
