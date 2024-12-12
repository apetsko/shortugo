package inmem

import (
	"errors"
)

type InMemStorage struct {
	data map[string]string
}

func New() *InMemStorage {
	return &InMemStorage{data: make(map[string]string)}
}

func (im *InMemStorage) Put(id, url string) (err error) {
	im.data[id] = url
	return nil
}

func (im *InMemStorage) Get(id string) (url string, err error) {
	if url, ok := im.data[id]; ok {
		return url, nil
	}
	return "", errors.New("URL not found")
}
