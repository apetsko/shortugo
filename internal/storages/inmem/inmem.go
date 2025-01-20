package inmem

import (
	"context"
	"fmt"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
)

type Storage struct {
	data map[string]string
}

func New() *Storage {
	return &Storage{data: make(map[string]string)}
}

func (im *Storage) Put(ctx context.Context, r models.URLRecord) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		im.data[r.ID] = r.URL
	}

	return nil
}

func (im *Storage) PutBatch(ctx context.Context, rr []models.URLRecord) (err error) {
	for _, r := range rr {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			im.data[r.ID] = r.URL
		}
	}
	return nil
}

func (im *Storage) Get(ctx context.Context, shortURL string) (url string, err error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if url, ok := im.data[shortURL]; ok {
			return url, nil
		}
	}
	return "", fmt.Errorf("URL not found: %s. %w", shortURL, shared.ErrNotFound)
}

func (im *Storage) Ping() error {
	return nil
}

func (im *Storage) Close() error {
	return nil
}
