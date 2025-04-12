// Package storages provides the initialization and management of different storage implementations.
// It includes logic for selecting the appropriate storage type and handling batch operations.
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

// Storage interface defines the methods that any storage implementation must provide.
type Storage interface {
	// DeleteUserURLs deletes multiple URLs associated with a user ID.
	DeleteUserURLs(ctx context.Context, IDs []string, userID string) (err error)
}

// Init initializes the appropriate storage based on the provided configuration.
func Init(databaseDSN, fileStoragePath string, logger *logging.Logger) (handlers.Storage, error) {
	switch {
	case databaseDSN != "":
		// Initialize PostgreSQL storage if databaseDSN is provided.
		s, err := postgres.New(databaseDSN, logger)
		if err != nil {
			return nil, err
		}
		logger.Info("Using database storages")
		return s, nil
	case fileStoragePath != "":
		// Initialize file storage if fileStoragePath is provided.
		s, err := infile.New(fileStoragePath)
		if err != nil {
			return nil, err
		}
		logger.Infof("Using file storage: %s", fileStoragePath)
		return s, nil
	default:
		// Initialize in-memory storage if no other storage configuration is provided.
		s := inmem.New()
		logger.Info("Using in-memory storages")
		return s, nil
	}
}

// StartBatchDeleteProcessor starts a background processor to handle batch delete requests.
func StartBatchDeleteProcessor(ctx context.Context, s Storage, input <-chan models.BatchDeleteRequest, logger *logging.Logger) {
	const (
		batchSize = 100             // Maximum number of requests to process in a single batch.
		timeout   = 2 * time.Second // Time interval to flush the batch if not full.
	)

	var batch []models.BatchDeleteRequest
	var mu sync.Mutex
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case req := <-input:
			// Add request to the batch.
			mu.Lock()
			batch = append(batch, req)
			if len(batch) >= batchSize {
				// Flush the batch if it reaches the batch size.
				flushBatch(ctx, s, &batch, logger)
			}
			mu.Unlock()
		case <-ticker.C:
			// Flush the batch at regular intervals.
			mu.Lock()
			if len(batch) > 0 {
				flushBatch(ctx, s, &batch, logger)
			}
			mu.Unlock()
		case <-ctx.Done():
			// Exit the processor when the context is done.
			return
		}
	}
}

// flushBatch processes and deletes URLs in the batch.
func flushBatch(ctx context.Context, s Storage, batch *[]models.BatchDeleteRequest, logger *logging.Logger) {
	if len(*batch) == 0 {
		return
	}

	defer func() { *batch = nil }()

	var wg sync.WaitGroup
	for _, b := range *batch {
		wg.Add(1)
		go func(userID string, ids []string) {
			defer wg.Done()
			// Delete URLs for the user.
			if err := s.DeleteUserURLs(ctx, ids, userID); err != nil {
				logger.Error(fmt.Errorf("error deleting URLs for user %s: %w", userID, err).Error())
			}
		}(b.UserID, b.Ids)
	}
	wg.Wait()
}
