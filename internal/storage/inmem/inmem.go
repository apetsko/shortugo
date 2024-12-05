package inmem

import (
	"errors"

	"github.com/apetsko/shortugo/internal/utils"
)

type InMemStorage struct {
	data map[string]string
}

func New() *InMemStorage {
	return &InMemStorage{data: make(map[string]string)}
}

func (im InMemStorage) Put(URL string) (ID string, err error) {
	if ID, err = utils.Generate(URL); err != nil {
		return "", err
	}
	im.data[ID] = URL
	return ID, nil
}

func (im InMemStorage) Get(ID string) (URL string, err error) {
	if URL, ok := im.data[ID]; ok {
		return URL, nil
	}
	return "", errors.New("URL not found")
}
