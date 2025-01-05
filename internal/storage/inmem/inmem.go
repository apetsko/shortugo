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

func (im *InMemStorage) Put(ID, URL string) (err error) {
	im.data[ID] = URL
	return nil
}

func (im *InMemStorage) Get(ID string) (URL string, err error) {
	if URL, ok := im.data[ID]; ok {
		return URL, nil
	}
	return "", errors.New("URL not found")
}
