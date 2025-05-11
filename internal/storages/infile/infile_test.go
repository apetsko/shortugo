package infile

import (
	"context"
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
		if err := store.Close(); err != nil {
			t.Errorf("failed to close storage: %v", err)
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	}
}

func TestStorage_PutAndGet(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	record := models.URLRecord{ID: "short123", URL: "http://example.com", UserID: "user1"}

	err := store.Put(ctx, record)
	require.NoError(t, err)

	url, err := store.Get(ctx, "short123")
	require.NoError(t, err)
	assert.Equal(t, "http://example.com", url)
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
	require.NoError(t, err)

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
	assert.ErrorIs(t, err, shared.ErrNotFound)
}

func TestStorage_Get_Deleted(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	record := models.URLRecord{ID: "del123", URL: "http://example.com", UserID: "user1", Deleted: true}
	require.NoError(t, store.Put(ctx, record))

	_, err := store.Get(ctx, "del123")
	assert.ErrorContains(t, err, "Gone")
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
	assert.Equal(t, "http://short.ly/short1", links[0].ID)
	assert.Equal(t, "http://short.ly/short2", links[1].ID)
}

func TestStorage_ListLinksByUserID_NotFound(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	_, err := store.ListLinksByUserID(ctx, "http://short.ly", "ghost")
	assert.ErrorIs(t, err, shared.ErrNotFound)
}

func TestStorage_DeleteUserURLs(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	records := []models.URLRecord{
		{ID: "short1", URL: "http://one.com", UserID: "user1"},
		{ID: "short2", URL: "http://two.com", UserID: "user1"},
	}
	require.NoError(t, store.PutBatch(ctx, records))

	err := store.DeleteUserURLs(ctx, []string{"short1"}, "user1")
	require.NoError(t, err)

	_, err = store.Get(ctx, "short1")
	assert.ErrorContains(t, err, "Gone")

	url, err := store.Get(ctx, "short2")
	require.NoError(t, err)
	assert.Equal(t, "http://two.com", url)
}

func TestStorage_Ping(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	assert.NoError(t, store.Ping())
}

func TestStorage_Stats(t *testing.T) {
	store, cleanup := setupTempStorage(t)
	defer cleanup()

	ctx := context.Background()
	records := []models.URLRecord{
		{ID: "a", URL: "http://a.com", UserID: "user1"},
		{ID: "b", URL: "http://b.com", UserID: "user2"},
		{ID: "c", URL: "http://c.com", UserID: "user1"},
	}

	require.NoError(t, store.PutBatch(ctx, records))

	stats, err := store.Stats(ctx)
	require.NoError(t, err)
	assert.Equal(t, 3, stats.Urls)
	assert.Equal(t, 2, stats.Users)
}

func TestCustomBool_JSON(t *testing.T) {
	var b CustomBool

	// Test Unmarshal
	require.NoError(t, b.UnmarshalJSON([]byte("1")))
	assert.True(t, bool(b))

	require.NoError(t, b.UnmarshalJSON([]byte("0")))
	assert.False(t, bool(b))

	require.NoError(t, b.UnmarshalJSON([]byte{}))
	assert.False(t, bool(b))

	// Test Marshal
	b = true
	out, err := b.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte("1"), out)

	b = false
	out, err = b.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte("0"), out)
}
