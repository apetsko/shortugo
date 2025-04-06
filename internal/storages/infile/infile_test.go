package infile

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTempStorage(t *testing.T) (*Storage, func()) {
	tmpFile, err := os.CreateTemp("", "test_storage")
	require.NoError(t, err, "failed to create temp file")

	store, err := New(tmpFile.Name())
	require.NoError(t, err, "failed to create storage")

	return store, func() {
		store.Close()
		os.Remove(tmpFile.Name())
	}
}

func TestStorage_PutAndGet(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	record := models.URLRecord{ID: "short123", URL: "http://example.com", UserID: "user1"}

	err := store.Put(ctx, record)
	require.NoError(t, err, "Put failed")

	url, err := store.Get(ctx, "short123")
	require.NoError(t, err, "Get failed")
	assert.Equal(t, "http://example.com", url, "URL mismatch")
}

func TestStorage_PutBatchAndGet(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	records := []models.URLRecord{
		{ID: "short1", URL: "http://one.com", UserID: "user1"},
		{ID: "short2", URL: "http://two.com", UserID: "user2"},
	}

	err := store.PutBatch(ctx, records)
	require.NoError(t, err, "PutBatch failed")

	url1, err := store.Get(ctx, "short1")
	require.NoError(t, err)
	assert.Equal(t, "http://one.com", url1)

	url2, err := store.Get(ctx, "short2")
	require.NoError(t, err)
	assert.Equal(t, "http://two.com", url2)
}

func TestStorage_Get_NotFound(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	_, err := store.Get(ctx, "nonexistent")
	assert.ErrorIs(t, err, shared.ErrNotFound, "expected ErrNotFound")
}

func TestStorage_ListLinksByUserID(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	records := []models.URLRecord{
		{ID: "short1", URL: "http://one.com", UserID: "user1"},
		{ID: "short2", URL: "http://two.com", UserID: "user1"},
		{ID: "short3", URL: "http://three.com", UserID: "user2"},
	}

	err := store.PutBatch(ctx, records)
	require.NoError(t, err)

	links, err := store.ListLinksByUserID(ctx, "http://short.ly", "user1")
	require.NoError(t, err)
	assert.Len(t, links, 2)

	expectedIDs := []string{"http://short.ly/short1", "http://short.ly/short2"}
	assert.Equal(t, expectedIDs[0], links[0].ID)
	assert.Equal(t, expectedIDs[1], links[1].ID)
}

func TestStorage_DeleteUserURLs(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	records := []models.URLRecord{
		{ID: "short1", URL: "http://one.com", UserID: "user1"},
		{ID: "short2", URL: "http://two.com", UserID: "user1"},
	}

	err := store.PutBatch(ctx, records)
	require.NoError(t, err)

	err = store.DeleteUserURLs(ctx, []string{"short1"}, "user1")
	require.NoError(t, err)

	str, err := store.Get(ctx, "short1")
	fmt.Println("str", str)
	assert.Error(t, err, "Expected error for deleted URL")
}

func TestStorage_Ping(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	err := store.Ping()
	assert.NoError(t, err, "Ping should always return nil")
}
