package postgres

import (
	"context"
	"fmt"
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
	ciConnString       = "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable" // Used in CI environment
)

var (
	logger, _ = logging.New(zapcore.DebugLevel)
	connStr   string
	isCI      = os.Getenv("CI") == "true"
)

// startTestDB starts a PostgreSQL database in a Docker container for local testing.
// If running in a CI environment, it uses the pre-configured CI connection string.
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

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to start PostgreSQL: %v\n", err)
		os.Exit(1)
	}

	waitForTestDB()
}

// stopTestDB stops the Docker container running the test database.
// This is only executed in a local environment, not in CI.
func stopTestDB() {
	if isCI {
		return
	}

	logger.Info("🛑 Stopping local test database...")
	cmd := exec.Command("docker", "stop", localContainerName)
	_ = cmd.Run() // Ignore errors as the container might not exist
}

// waitForTestDB waits for the PostgreSQL database to become ready.
// It retries the connection until a timeout is reached.
func waitForTestDB() {
	logger.Info("⏳ Waiting for PostgreSQL to become ready...")

	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			logger.Info("❌ Timeout while waiting for PostgreSQL")
			os.Exit(1)
		case <-tick:
			storage, err := New(connStr, logger)
			if err == nil {
				if closeErr := storage.Close(); closeErr != nil {
					logger.Errorf("Failed to create & close PostgreSQL storage: %s. %s", closeErr, err)
					return
				}
				logger.Info("✅ PostgreSQL is ready")
				return
			}
		}
	}
}

// setupTestStorage creates a new Storage instance for testing.
// It ensures the storage is properly closed after the test.
func setupTestStorage(t *testing.T) *Storage {
	logger.Info(fmt.Sprintf("LOG setupTestStorage: connecting to %q", connStr))
	storage, err := New(connStr, logger)
	require.NoError(t, err)
	t.Cleanup(func() { _ = storage.Close() })
	return storage
}

// TestMain manages the lifecycle of tests.
// It starts the test database if running locally and stops it after tests are completed.
func TestMain(m *testing.M) {
	isCI := os.Getenv("CI") == "true"
	logger.Infof("CI.env in test: %q", os.Getenv("CI"))

	if isCI {
		connStr = ciConnString
		logger.Info("🔄 Running in CI. Using CI Postgres container...", "connStr", connStr)
	} else {
		connStr = localConnString
		logger.Info("🔄 Starting local test DB...", "connStr", connStr)
		startTestDB()
		time.Sleep(5 * time.Second)
	}

	exitCode := m.Run()

	if !isCI {
		stopTestDB()
	}

	os.Exit(exitCode)
}

// TestStorage_PutGet tests the Put and Get methods of the Storage.
// It verifies that a URLRecord can be stored and retrieved successfully.
func TestStorage_PutGet(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	record := models.URLRecord{
		ID:     "test-id",
		URL:    "https://example.com",
		UserID: "user-123",
	}

	err := storage.Put(ctx, record)
	require.NoError(t, err)

	gotURL, err := storage.Get(ctx, "test-id")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", gotURL)
}

// TestStorage_DeleteUserURLs tests the DeleteUserURLs method of the Storage.
// It ensures that URLs associated with a user can be marked as deleted.
func TestStorage_DeleteUserURLs(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	record := models.URLRecord{
		ID:     "del-id",
		URL:    "https://delete.com",
		UserID: "user-456",
	}

	err := storage.Put(ctx, record)
	require.NoError(t, err)

	err = storage.DeleteUserURLs(ctx, []string{"del-id"}, "user-456")
	require.NoError(t, err)

	_, err = storage.Get(ctx, "del-id")
	assert.ErrorIs(t, err, shared.ErrGone)
}

// TestStorage_PutBatch tests the PutBatch method of the Storage.
// It verifies that multiple URLRecords can be stored and retrieved successfully.
func TestStorage_PutBatch(t *testing.T) {
	storage := setupTestStorage(t)
	ctx := context.Background()

	records := []models.URLRecord{
		{ID: "batch1", URL: "https://batch1.com", UserID: "user-789"},
		{ID: "batch2", URL: "https://batch2.com", UserID: "user-789"},
	}

	err := storage.PutBatch(ctx, records)
	require.NoError(t, err)

	gotURL1, err := storage.Get(ctx, "batch1")
	require.NoError(t, err)
	assert.Equal(t, "https://batch1.com", gotURL1)

	gotURL2, err := storage.Get(ctx, "batch2")
	require.NoError(t, err)
	assert.Equal(t, "https://batch2.com", gotURL2)
}
