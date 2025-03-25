package storages

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/infile"
	"github.com/apetsko/shortugo/internal/storages/inmem"
	"github.com/apetsko/shortugo/internal/storages/postgres"
)

type Storage interface {
	DeleteUserURLs(ctx context.Context, IDs []string, userID string) (err error)
}

func Init(databaseDSN, fileStoragePath string, logger *logging.ZapLogger) (handlers.Storage, error) {
	switch {
	case databaseDSN != "":
		s, err := postgres.New(databaseDSN)
		if err != nil {
			return nil, err
		}
		logger.Info("Using database storages")
		return s, nil
	case fileStoragePath != "":
		s, err := infile.New(fileStoragePath)
		if err != nil {
			return nil, err
		}
		logger.Infof("Using file storage: %s", fileStoragePath)
		return s, nil
	default:
		s := inmem.New()
		logger.Info("Using in-memory storages")
		return s, nil
	}
}

func StartBatchDeleteProcessor(ctx context.Context, s Storage, input <-chan models.BatchDeleteRequest, logger *logging.ZapLogger) {
	const (
		batchSize = 100
		timeout   = 2 * time.Second
	)

	var batch []models.BatchDeleteRequest
	var mu sync.Mutex
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case req := <-input:
			mu.Lock()
			batch = append(batch, req)
			if len(batch) >= batchSize {
				flushBatch(ctx, s, &batch, logger)
			}
			mu.Unlock()
		case <-ticker.C:
			mu.Lock()
			if len(batch) > 0 {
				flushBatch(ctx, s, &batch, logger)
			}
			mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func flushBatch(ctx context.Context, s Storage, batch *[]models.BatchDeleteRequest, logger *logging.ZapLogger) {
	if len(*batch) == 0 {
		return
	}

	defer func() { *batch = nil }()

	var wg sync.WaitGroup
	for _, b := range *batch {
		wg.Add(1)
		go func(userID string, ids []string) {
			defer wg.Done()
			if err := s.DeleteUserURLs(ctx, ids, userID); err != nil {
				logger.Error(fmt.Errorf("error deleting URLs for user %s: %w", userID, err).Error())
			}
		}(b.UserID, b.Ids)
	}
	wg.Wait()
}
