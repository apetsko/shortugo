package infile

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
)

const FilePermUserRWGroupROthersR = 0644

type CustomBool bool

func (b *CustomBool) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*b = false
		return nil
	}

	if data[0] == '1' {
		*b = true
	} else {
		*b = false
	}
	return nil
}

func (b CustomBool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

type Storage struct {
	file    *os.File
	encoder *json.Encoder
	mu      sync.Mutex
}

func New(filename string) (*Storage, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, FilePermUserRWGroupROthersR)
	if err != nil {
		return nil, err
	}

	return &Storage{
		file:    f,
		encoder: json.NewEncoder(f),
	}, nil
}

func (f *Storage) Close() error {
	return f.file.Close()
}

func (f *Storage) Put(ctx context.Context, r models.URLRecord) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.encoder.Encode(r); err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := f.file.Sync(); err != nil {
		return fmt.Errorf("error sync file: %w", err)
	}

	return nil
}

func (f *Storage) PutBatch(ctx context.Context, rr []models.URLRecord) (err error) {
	for _, r := range rr {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := f.encoder.Encode(r); err != nil {
			return err
		}

		if err := f.file.Sync(); err != nil {
			return fmt.Errorf("error sync file: %w", err)
		}
	}
	return nil
}

func (f *Storage) Get(ctx context.Context, shortURL string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if _, err := f.file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("error setting file seek: %w", err)
	}

	r := new(models.URLRecord)
	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		err := json.Unmarshal(data, &r)
		if err != nil {
			return "", fmt.Errorf("failed unmarshal: %w", err)
		}

		if r.ID == shortURL {
			if r.Deleted {
				return "", errors.New(http.StatusText(http.StatusGone))
			}
			return r.URL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return "", fmt.Errorf("URL not found: %s. %w", shortURL, shared.ErrNotFound)
}

func (f *Storage) ListLinksByUserID(ctx context.Context, userID, baseURL string) ([]models.URLRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if _, err := f.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error setting file seek: %w", err)
	}

	rr := make([]models.URLRecord, 0)
	r := new(models.URLRecord)
	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		if err := json.Unmarshal(data, r); err != nil {
			return nil, err
		}

		if r.UserID == userID && !r.Deleted {
			r.ID = fmt.Sprintf("%s/%s", baseURL, r.ID)
			rr = append(rr, *r)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(rr) == 0 {
		return rr, fmt.Errorf("URLs not found for UserID: %s. %w", userID, shared.ErrNotFound)
	}
	return rr, nil
}

func (f *Storage) DeleteUserURLs(ctx context.Context, ids []string, userID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := f.file.Seek(0, 0); err != nil {
		return fmt.Errorf("error setting file seek: %w", err)
	}

	var offset int64

	r := new(models.URLRecord)
	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		if err := json.Unmarshal(data, r); err != nil {
			return err
		}

		if r.UserID == userID && slices.Contains(ids, r.ID) && !r.Deleted {
			res := strings.Replace(string(data), `"deleted":false`, `"deleted":true `, 1)
			if _, err := f.file.Seek(offset, 0); err != nil {
				return fmt.Errorf("error setting file seek: %w", err)
			}
			if _, err := f.file.WriteString(res); err != nil {
				return fmt.Errorf("error writing updated record to file: %w", err)
			}
			offset, _ = f.file.Seek(0, 1)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	return nil
}

func (f *Storage) Ping() error {
	return nil
}
