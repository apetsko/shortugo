package postgres

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"testing"
// 	"time"
//
// 	"github.com/apetsko/shortugo/internal/logging"
// 	"github.com/apetsko/shortugo/internal/models"
// 	"github.com/apetsko/shortugo/internal/storages/shared"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/zap/zapcore"
// )
//
// const (
// 	localContainerName = "test_postgres_container"
// 	localConnString    = "postgres://testuser:testpass@localhost:54321/testdb?sslmode=disable"
// 	ciConnString       = "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
// )
//
// var (
// 	logger, _ = logging.New(zapcore.DebugLevel)
// 	connStr   string
// 	isCI      = os.Getenv("CI") == "true"
// )
//
// // startTestDB запускает БД в Docker, если тесты выполняются локально
// func startTestDB() {
// 	if isCI {
// 		connStr = ciConnString
// 		logger.Info("🔄 Тесты запущены в CI. Используем встроенный PostgreSQL...")
// 		return
// 	}
//
// 	connStr = localConnString
//
// 	logger.Info("Используем строку подключения:", connStr)
//
// 	logger.Info("🔄 Запускаем тестовую базу данных в Docker...")
//
// 	cmd := exec.Command("docker", "run", "--rm", "-d",
// 		"--name", localContainerName,
// 		"-e", "POSTGRES_DB=testdb",
// 		"-e", "POSTGRES_USER=testuser",
// 		"-e", "POSTGRES_PASSWORD=testpass",
// 		"-p", "54321:5432",
// 		"postgres:15")
//
// 	if err := cmd.Run(); err != nil {
// 		fmt.Printf("❌ Ошибка запуска PostgreSQL: %v\n", err)
// 		os.Exit(1)
// 	}
//
// 	waitForDB()
// }
//
// // stopTestDB останавливает контейнер после тестов (если они не в CI)
// func stopTestDB() {
// 	if isCI {
// 		return
// 	}
//
// 	logger.Info("🛑 Останавливаем локальную тестовую базу данных...")
// 	cmd := exec.Command("docker", "stop", localContainerName)
// 	_ = cmd.Run() // Ошибки игнорируем, т.к. контейнера может не быть
// }
//
// // waitForDB ожидает готовность PostgreSQL
// func waitForDB() {
// 	logger.Info("⏳ Ожидаем готовность PostgreSQL...")
//
// 	timeout := time.After(30 * time.Second)
// 	tick := time.Tick(1 * time.Second)
//
// 	for {
// 		select {
// 		case <-timeout:
// 			logger.Info("❌ Тайм-аут ожидания PostgreSQL")
// 			os.Exit(1)
// 		case <-tick:
// 			storage, err := New(connStr)
// 			if err == nil {
// 				storage.Close()
// 				logger.Info("✅ PostgreSQL готов к работе")
// 				return
// 			}
// 		}
// 	}
// }
//
// // setupTestStorage создает хранилище для тестов
// func setupTestStorage(t *testing.T) *Storage {
// 	storage, err := New(connStr)
// 	require.NoError(t, err)
// 	t.Cleanup(func() { _ = storage.Close() })
// 	return storage
// }
//
// // TestMain управляет жизненным циклом тестов
// func TestMain(m *testing.M) {
// 	if !isCI {
// 		startTestDB()
// 	}
//
// 	exitCode := m.Run()
//
// 	if !isCI {
// 		stopTestDB()
// 	}
//
// 	os.Exit(exitCode)
// }
//
// // Тесты
//
// func TestStorage_PutGet(t *testing.T) {
// 	storage := setupTestStorage(t)
// 	ctx := context.Background()
//
// 	record := models.URLRecord{
// 		ID:     "test-id",
// 		URL:    "https://example.com",
// 		UserID: "user-123",
// 	}
//
// 	err := storage.Put(ctx, record)
// 	require.NoError(t, err)
//
// 	gotURL, err := storage.Get(ctx, "test-id")
// 	require.NoError(t, err)
// 	assert.Equal(t, "https://example.com", gotURL)
// }
//
// func TestStorage_DeleteUserURLs(t *testing.T) {
// 	storage := setupTestStorage(t)
// 	ctx := context.Background()
//
// 	record := models.URLRecord{
// 		ID:     "del-id",
// 		URL:    "https://delete.com",
// 		UserID: "user-456",
// 	}
//
// 	err := storage.Put(ctx, record)
// 	require.NoError(t, err)
//
// 	err = storage.DeleteUserURLs(ctx, []string{"del-id"}, "user-456")
// 	require.NoError(t, err)
//
// 	_, err = storage.Get(ctx, "del-id")
// 	assert.ErrorIs(t, err, shared.ErrGone)
// }
//
// func TestStorage_PutBatch(t *testing.T) {
// 	storage := setupTestStorage(t)
// 	ctx := context.Background()
//
// 	records := []models.URLRecord{
// 		{ID: "batch1", URL: "https://batch1.com", UserID: "user-789"},
// 		{ID: "batch2", URL: "https://batch2.com", UserID: "user-789"},
// 	}
//
// 	err := storage.PutBatch(ctx, records)
// 	require.NoError(t, err)
//
// 	gotURL1, err := storage.Get(ctx, "batch1")
// 	require.NoError(t, err)
// 	assert.Equal(t, "https://batch1.com", gotURL1)
//
// 	gotURL2, err := storage.Get(ctx, "batch2")
// 	require.NoError(t, err)
// 	assert.Equal(t, "https://batch2.com", gotURL2)
// }
