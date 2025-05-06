// Package postgres provides a PostgreSQL-backed storage implementation for the application.
// It includes methods for storing, retrieving, and managing URL records in a PostgreSQL database.
package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Storage represents a PostgreSQL-backed storage.
type Storage struct {
	pool *pgxpool.Pool
}

// applyMigrations applies database migrations using goose.
func applyMigrations(conn string, logger *logging.Logger) error {
	goose.SetBaseFS(migrations)
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return fmt.Errorf("goose: failed to open DB: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()

	err = goose.Up(db, "migrations")
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// New creates a new Storage instance and applies migrations.
func New(conn string, logger *logging.Logger) (*Storage, error) {
	// Подождем, пока БД станет доступна (в CI может занять пару секунд)
	if err := waitForDB(conn, logger); err != nil {
		return nil, fmt.Errorf("database not ready: %w", err)
	}

	// Применяем миграции
	if err := applyMigrations(conn, logger); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Storage{pool: pool}, nil
}

// Close closes the database connection pool.
func (p *Storage) Close() error {
	p.pool.Close()
	return nil
}

// waitForDB пытается подключиться к БД несколько раз с интервалами
func waitForDB(conn string, logger *logging.Logger) error {
	var lastErr error
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for i := 0; i < 10; i++ {
		pool, err := pgxpool.New(ctx, conn)
		if err == nil {
			pool.Close()
			logger.Info("✅ Database is ready")
			return nil
		}

		lastErr = err
		logger.Infof("⏳ Waiting for database... (%d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("database not ready after retries: %w", lastErr)
}

// Put stores a URLRecord in the database.
func (p *Storage) Put(ctx context.Context, r models.URLRecord) error {
	const insert = `
			INSERT INTO urls (id, url, user_id, date)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id)
			DO UPDATE SET date = EXCLUDED.date;`

	_, err := p.pool.Exec(ctx, insert, r.ID, r.URL, r.UserID, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to insert URL: %w", err)
	}

	return nil
}

// PutBatch stores multiple URLRecords in the database.
func (p *Storage) PutBatch(ctx context.Context, rr []models.URLRecord) error {
	const insertBatch = `
			INSERT INTO urls (id, url, user_id, date)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id)
			DO UPDATE SET date = EXCLUDED.date;`

	batch := new(pgx.Batch)
	for _, r := range rr {
		batch.Queue(insertBatch, r.ID, r.URL, r.UserID, time.Now().Format(time.RFC3339))
	}
	br := p.pool.SendBatch(ctx, batch)
	defer func() {
		_ = br.Close()
	}()

	if _, err := br.Exec(); err != nil {
		return fmt.Errorf("failed to batch insert: %w", err)
	}

	return nil
}

// Get retrieves the original URL for a given short URL.
func (p *Storage) Get(ctx context.Context, id string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	deleted := false
	const query = "SELECT url, deleted FROM urls WHERE id = $1"

	var url string
	err := p.pool.QueryRow(ctx, query, id).Scan(&url, &deleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("URL not found: %s. %w", id, shared.ErrNotFound)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return "", err
	}

	if deleted {
		return "", shared.ErrGone
	}

	return url, nil
}

// ListLinksByUserID lists all URLs associated with a user ID.
func (p *Storage) ListLinksByUserID(ctx context.Context, baseURL, userID string) ([]models.URLRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const query = "SELECT id, url, user_id FROM urls WHERE user_id = $1 AND deleted = FALSE"

	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	rr := make([]models.URLRecord, 0)
	for rows.Next() {
		var record models.URLRecord
		if err := rows.Scan(&record.ID, &record.URL, &record.UserID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		record.ID = baseURL + "/" + record.ID
		rr = append(rr, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(rr) == 0 {
		return nil, fmt.Errorf("urls not found for userID: %s. %w", userID, shared.ErrNotFound)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return rr, nil
}

// DeleteUserURLs deletes multiple URLs associated with a user ID.
func (p *Storage) DeleteUserURLs(ctx context.Context, ids []string, userID string) error {
	const setDeleteBatch = `
				UPDATE urls
				SET deleted = true
				WHERE id = ANY($1::text[]) AND user_id = $2 AND deleted = FALSE;`

	batch := new(pgx.Batch)
	batch.Queue(setDeleteBatch, ids, userID)

	br := p.pool.SendBatch(ctx, batch)
	defer func() {
		_ = br.Close()
	}()
	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("failed to batch delete user urls: %w", err)
	}

	return nil
}

// Stats retrieves count stats: urls and users.
func (p *Storage) Stats(ctx context.Context) (*models.Stats, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const query = `
		SELECT COUNT(*) AS total_links, COUNT(DISTINCT user_id) AS unique_users
		FROM urls;
	`

	var totalLinks int
	var uniqueUsers int

	err := p.pool.QueryRow(ctx, query).Scan(&totalLinks, &uniqueUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return &models.Stats{
		Urls:  totalLinks,
		Users: uniqueUsers,
	}, nil
}

// Ping checks the database connection health.
func (p *Storage) Ping() error {
	return p.pool.Ping(context.Background())
}
