package postgres

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

const (
	localContainerName = "test_postgres_container"
	localConnString    = "postgres://testuser:testpass@localhost:54321/testdb?sslmode=disable"
	ciConnString       = "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
)

var (
	logger, _ = logging.New(zapcore.DebugLevel)
	connStr   string
	isCI      = os.Getenv("CI") == "true"
)

func startTestDB() {
	if isCI {
		connStr = ciConnString
		logger.Info("🔄 Tests are running in CI. Using built-in PostgreSQL...")
		return
	}

	connStr = localConnString
	logger.Info("🔄 Starting test database in Docker...")

	cmd := exec.Command("docker", "run", "--rm", "-d",
		"--name", localContainerName,
		"-e", "POSTGRES_DB=testdb",
		"-e", "POSTGRES_USER=testuser",
		"-e", "POSTGRES_PASSWORD=testpass",
		"-p", "54321:5432",
		"postgres:17")
	require.NoError(nil, cmd.Run())
	waitForTestDB()
}

func stopTestDB() {
	if !isCI {
		logger.Info("🛑 Stopping local test database...")
		_ = exec.Command("docker", "stop", localContainerName).Run()
	}
}

func waitForTestDB() {
	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-timeout:
			logger.Error("❌ Timeout while waiting for PostgreSQL")
			os.Exit(1)
		case <-tick:
			storage, err := New(connStr, logger)
			if err == nil {
				_ = storage.Close()
				logger.Info("✅ PostgreSQL is ready")
				return
			}
		}
	}
}

func setupTestStorage(t *testing.T) *Storage {
	storage, err := New(connStr, logger)
	require.NoError(t, err)
	t.Cleanup(func() { _ = storage.Close() })
	return storage
}

func TestMain(m *testing.M) {
	if isCI {
		connStr = ciConnString
	} else {
		connStr = localConnString
		startTestDB()
		time.Sleep(5 * time.Second)
	}
	code := m.Run()
	stopTestDB()
	os.Exit(code)
}

func TestStorage_PutGet(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	rec := models.URLRecord{ID: "id-get", URL: "https://example.com", UserID: "user1"}
	require.NoError(t, storage.Put(ctx, rec))

	url, err := storage.Get(ctx, "id-get")
	require.NoError(t, err)
	assert.Equal(t, rec.URL, url)
}

func TestStorage_Get_NotFound(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	_, err := storage.Get(ctx, "unknown-id")
	assert.ErrorIs(t, err, shared.ErrNotFound)
}

func TestStorage_Get_Deleted(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	rec := models.URLRecord{ID: "id-del", URL: "https://del.com", UserID: "user-del"}
	require.NoError(t, storage.Put(ctx, rec))
	require.NoError(t, storage.DeleteUserURLs(ctx, []string{"id-del"}, "user-del"))

	_, err := storage.Get(ctx, "id-del")
	assert.ErrorIs(t, err, shared.ErrGone)
}

func TestStorage_PutBatch(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	records := []models.URLRecord{
		{ID: "b1", URL: "https://b1.com", UserID: "user"},
		{ID: "b2", URL: "https://b2.com", UserID: "user"},
	}
	require.NoError(t, storage.PutBatch(ctx, records))

	for _, r := range records {
		url, err := storage.Get(ctx, r.ID)
		require.NoError(t, err)
		assert.Equal(t, r.URL, url)
	}
}

func TestStorage_DeleteUserURLs(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	rec := models.URLRecord{ID: "to-delete", URL: "https://del.com", UserID: "u-del"}
	require.NoError(t, storage.Put(ctx, rec))
	require.NoError(t, storage.DeleteUserURLs(ctx, []string{"to-delete"}, "u-del"))

	_, err := storage.Get(ctx, "to-delete")
	assert.ErrorIs(t, err, shared.ErrGone)
}

func TestStorage_ListLinksByUserID(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	records := []models.URLRecord{
		{ID: "l1", URL: "http://1.com", UserID: "user-l"},
		{ID: "l2", URL: "http://2.com", UserID: "user-l"},
	}
	require.NoError(t, storage.PutBatch(ctx, records))

	links, err := storage.ListLinksByUserID(ctx, "http://short", "user-l")
	require.NoError(t, err)
	assert.Len(t, links, 2)
	assert.Contains(t, links[0].ID, "http://short/")
}

func TestStorage_ListLinksByUserID_NotFound(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	_, err := storage.ListLinksByUserID(ctx, "http://short", "ghost")
	assert.ErrorIs(t, err, shared.ErrNotFound)
}

func TestStorage_Stats(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	records := []models.URLRecord{
		{ID: "s1", URL: "https://a.com", UserID: "user1"},
		{ID: "s2", URL: "https://b.com", UserID: "user2"},
		{ID: "s3", URL: "https://c.com", UserID: "user1"},
	}
	require.NoError(t, storage.PutBatch(ctx, records))

	stats, err := storage.Stats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.Urls, 3)
	assert.GreaterOrEqual(t, stats.Users, 2)
}

func TestStorage_Ping(t *testing.T) {
	storage := setupTestStorage(t)
	assert.NoError(t, storage.Ping())
}

func TestStorage_CtxCancelled(t *testing.T) {
	storage := setupTestStorage(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := storage.Get(ctx, "any")
	assert.ErrorIs(t, err, context.Canceled)
}
