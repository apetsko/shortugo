package storage

import (
	"errors"

	"github.com/apetsko/shortugo/internal/utils"
)

type Storage interface {
	Put(string) (string, error)
	Get(string) (string, error)
}

type inMem struct {
	m map[string]string
}

func NewInMem() *inMem {
	im := &inMem{map[string]string{}}
	return im
}

func (im inMem) Put(URL string) (ID string, err error) {
	if ID, err = utils.Generate(URL); err != nil {
		return "", err
	}
	im.m[ID] = URL
	return ID, nil
}

func (im inMem) Get(ID string) (URL string, err error) {
	if URL, ok := im.m[ID]; ok {
		return URL, nil
	}
	return "", errors.New("URL not found")
}
