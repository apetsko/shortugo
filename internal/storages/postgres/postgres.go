package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Storage struct {
	pool *pgxpool.Pool
}

func applyMigrations(conn string) error {
	goose.SetBaseFS(embedMigrations)
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return fmt.Errorf("goose: failed to open DB: %w", err)
	}
	defer db.Close()

	err = goose.Up(db, "migrations")
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func New(conn string) (*Storage, error) {
	if err := applyMigrations(conn); err != nil {
		return nil, err
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Storage{pool: pool}, nil
}

func (p *Storage) Close() error {
	p.pool.Close()
	return nil
}

func (p *Storage) Put(ctx context.Context, r models.URLRecord) error {
	const insert = `
	INSERT INTO urls (id, url, date)
	VALUES ($1, $2, $3)
	ON CONFLICT (id)
	DO UPDATE SET date = EXCLUDED.date;`

	_, err := p.pool.Exec(ctx, insert, r.ID, r.URL, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to insert URL: %w", err)
	}

	return nil
}

func (p *Storage) PutBatch(ctx context.Context, rr []models.URLRecord) error {
	const insertBatch = `
	INSERT INTO urls (id, url, date)
	VALUES ($1, $2, $3)
	ON CONFLICT (id)
	DO UPDATE SET date = EXCLUDED.date;`

	batch := new(pgx.Batch)
	for _, r := range rr {
		batch.Queue(insertBatch, r.ID, r.URL, time.Now().Format(time.RFC3339))
	}
	br := p.pool.SendBatch(ctx, batch)
	defer br.Close()

	if _, err := br.Exec(); err != nil {
		return fmt.Errorf("failed to batch insert: %w", err)
	}

	return nil
}

func (p *Storage) Get(ctx context.Context, id string) (url string, err error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	const query = "SELECT url FROM urls WHERE id=$1"

	err = p.pool.QueryRow(ctx, query, id).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("URL not found: %s. %w", id, shared.ErrNotFound)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return "", err
	}

	return url, nil
}

func (p *Storage) GetByUserID(ctx context.Context, userID string) (rr []models.URLRecord, err error) {

	return nil, nil
}

func (p *Storage) Ping() error {
	return p.pool.Ping(context.Background())
}
