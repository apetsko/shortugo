package inmem

import (
	"context"
	"fmt"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
)

type Storage struct {
	byID     map[string]models.URLRecord
	byUserID map[string][]models.URLRecord
}

func New() *Storage {
	return &Storage{
		byID:     make(map[string]models.URLRecord),
		byUserID: make(map[string][]models.URLRecord),
	}
}

func (im *Storage) Put(ctx context.Context, r models.URLRecord) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		im.byID[r.ID] = r
		im.byUserID[r.UserID] = append(im.byUserID[r.UserID], r)
	}
	return nil
}

func (im *Storage) PutBatch(ctx context.Context, rr []models.URLRecord) (err error) {
	for _, r := range rr {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			im.byID[r.ID] = r
			im.byUserID[r.UserID] = append(im.byUserID[r.UserID], r)
		}
	}
	return nil
}

func (im *Storage) Get(ctx context.Context, shortURL string) (url string, err error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if rec, ok := im.byID[shortURL]; ok {
			if rec.Deleted {
				return "", shared.ErrGone
			}
			return rec.URL, nil
		}
	}
	return "", fmt.Errorf("URL not found: %s. %w", shortURL, shared.ErrNotFound)
}

func (im *Storage) ListLinksByUserID(ctx context.Context, baseURL, userID string) (rr []models.URLRecord, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if recs, ok := im.byUserID[userID]; ok {
			for i := range recs {
				if !recs[i].Deleted {
					recs[i].ID = fmt.Sprintf("%s/%s", baseURL, recs[i].ID)
				}
			}
			return recs, nil
		}
	}
	return nil, fmt.Errorf("URLs not found for UserID: %s. %w", userID, shared.ErrNotFound)
}

func (im *Storage) DeleteUserURLs(ctx context.Context, ids []string, userID string) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		for _, id := range ids {
			if rec, ok := im.byID[id]; ok && !rec.Deleted {
				if rec.UserID == userID {
					rec.Deleted = true
					im.byID[id] = rec
				}
			}
		}
		return nil
	}
}

func (im *Storage) Ping() error {
	return nil
}

func (im *Storage) Close() error {
	return nil
}
