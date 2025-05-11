package storages_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages"
	"github.com/apetsko/shortugo/internal/storages/infile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func setupLogger(t *testing.T) *logging.Logger {
	logger, err := logging.New(zapcore.ErrorLevel)
	require.NoError(t, err)
	return logger
}

func TestInit_InMemory(t *testing.T) {
	logger := setupLogger(t)
	store, err := storages.Init("", "", logger)
	require.NoError(t, err)
	assert.NotNil(t, store)
}

func TestInit_FileStorage(t *testing.T) {
	logger := setupLogger(t)

	tmp, err := os.CreateTemp("", "storage-*.json")
	require.NoError(t, err)
	defer func() {
		err = os.Remove(tmp.Name())
		require.NoError(t, err)
	}()
	store, err := storages.Init("", tmp.Name(), logger)
	require.NoError(t, err)
	assert.NotNil(t, store)

	_, ok := store.(*infile.Storage)
	assert.True(t, ok)
}

func TestInit_InvalidFilePath(t *testing.T) {
	logger := setupLogger(t)
	_, err := storages.Init("", "/invalid/path/storage.json", logger)
	require.Error(t, err)
}

func TestInit_PostgresFails(t *testing.T) {
	logger := setupLogger(t)
	_, err := storages.Init("invalid-dsn", "", logger)
	require.Error(t, err)
}

type mockStorage struct {
	Deleted [][]string
}

func (m *mockStorage) DeleteUserURLs(ctx context.Context, IDs []string, userID string) error {
	m.Deleted = append(m.Deleted, IDs)
	return nil
}

func TestStartBatchDeleteProcessor_ImmediateFlushOnLimit(t *testing.T) {
	logger := setupLogger(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock := &mockStorage{}
	ch := make(chan models.BatchDeleteRequest, 110)

	go storages.StartBatchDeleteProcessor(ctx, mock, ch, logger)

	// Засунем >100 записей — должен быть вызов сразу
	for i := 0; i < 101; i++ {
		ch <- models.BatchDeleteRequest{
			UserID: "user",
			Ids:    []string{string(rune(65 + (i % 26)))}, // A-Z по кругу
		}
	}

	time.Sleep(3 * time.Second)
	cancel()

	require.GreaterOrEqual(t, len(mock.Deleted), 1)
}

func TestStartBatchDeleteProcessor_GracefulShutdown(t *testing.T) {
	logger := setupLogger(t)
	ctx, cancel := context.WithCancel(context.Background())

	mock := &mockStorage{}
	ch := make(chan models.BatchDeleteRequest)

	go storages.StartBatchDeleteProcessor(ctx, mock, ch, logger)

	ch <- models.BatchDeleteRequest{
		UserID: "shutdown-test",
		Ids:    []string{"x1", "x2"},
	}

	// Завершаем до таймера — должен обработать в `ctx.Done` ветке
	cancel()
	time.Sleep(500 * time.Millisecond)

	require.Len(t, mock.Deleted, 1)
	assert.ElementsMatch(t, []string{"x1", "x2"}, mock.Deleted[0])
}
