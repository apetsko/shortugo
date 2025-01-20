package inmem

import (
	"fmt"

	"github.com/apetsko/shortugo/internal/storage/shared"
)

type InMemStorage struct {
	data map[string]string
}

func New() *InMemStorage {
	return &InMemStorage{data: make(map[string]string)}
}

func (im *InMemStorage) Put(shortURL, url string) (err error) {
	im.data[shortURL] = url
	return nil
}

func (im *InMemStorage) Get(shortURL string) (url string, err error) {
	if url, ok := im.data[shortURL]; ok {
		return url, nil
	}
	return "", fmt.Errorf("URL not found: %s. %w", shortURL, shared.ErrNotFound)
}

func (im *InMemStorage) Ping() error {
	return nil
}

func (im *InMemStorage) Close() error {
	return nil
}
