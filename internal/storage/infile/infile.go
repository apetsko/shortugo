package infile

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/apetsko/shortugo/internal/storage"
)

const FilePermUserRWGroupROthersR = 0644

type InFileStorage struct {
	file    *os.File
	encoder *json.Encoder
	scanner *bufio.Scanner
	uuid    func() string
}

type Link struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func New(filename string) (*InFileStorage, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, FilePermUserRWGroupROthersR)
	if err != nil {
		return nil, err
	}

	return &InFileStorage{
		file:    f,
		encoder: json.NewEncoder(f),
		uuid:    nextUUID(),
		scanner: bufio.NewScanner(f),
	}, nil
}

func (f *InFileStorage) Close() error {
	return f.file.Close()
}

func (f *InFileStorage) Put(id, url string) (err error) {
	l := Link{
		UUID:        f.uuid(),
		ShortURL:    id,
		OriginalURL: url,
	}

	if err := f.encoder.Encode(l); err != nil {
		return err
	}

	if err := f.file.Sync(); err != nil {
		return fmt.Errorf("error sync file: %w", err)
	}

	return nil
}

func (f *InFileStorage) Get(id string) (string, error) {
	if _, err := f.file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("error setting file seek: %w", err)
	}
	f.scanner = bufio.NewScanner(f.file)
	for f.scanner.Scan() {
		data := f.scanner.Bytes()

		l := &Link{}
		err := json.Unmarshal(data, l)
		if err != nil {
			return "", err
		}

		if l.ShortURL == id {
			return l.OriginalURL, nil
		}
	}

	if err := f.scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return "", fmt.Errorf("URL not found: %s. %w", id, storage.ErrNotFound)
}

func nextUUID() func() string {
	count := 0
	return func() string {
		count++
		return fmt.Sprintf("%d", count)
	}
}
