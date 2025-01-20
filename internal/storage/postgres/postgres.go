package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storage/shared"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func New(connString string) (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
		id TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		date DATE NOT NULL
	);`

	_, err = conn.Exec(context.Background(), createTableQuery)
	if err != nil {
		err := conn.Close(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to close db connection: %w", err)
		}
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Storage{conn: conn}, nil
}

func (p *Storage) Close() error {
	return p.conn.Close(context.Background())
}

func (p *Storage) Put(ctx context.Context, r models.URLRecord) error {
	tx, err := p.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("transaction rollback failed: %w, original error: %v", rollbackErr, err)
			}
		}
	}()

	insertQuery := `
	INSERT INTO urls (id, url, date)
	VALUES ($1, $2, $3)
	ON CONFLICT (id)
	DO UPDATE SET date = EXCLUDED.date;`

	_, err = tx.Exec(context.Background(), insertQuery, r.ID, r.URL, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to insert URL: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (p *Storage) PutBatch(ctx context.Context, rr []models.URLRecord) error {
	tx, err := p.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("transaction rollback failed: %w, original error: %v", rollbackErr, err)
			}
		}
	}()

	insertQuery := `
	INSERT INTO urls (id, url, date)
	VALUES ($1, $2, $3)
	ON CONFLICT (id)
	DO UPDATE SET date = EXCLUDED.date;`

	_, err = tx.Prepare(ctx, "PutBatch", insertQuery)
	if err != nil {
		return err
	}

	for _, r := range rr {
		_, err := tx.Exec(ctx, "PutBatch", r.ID, r.URL, time.Now().Format(time.RFC3339))
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (p *Storage) Get(ctx context.Context, id string) (url string, err error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	err = p.conn.QueryRow(context.Background(), "SELECT url FROM urls WHERE id=$1", id).Scan(&url)
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

func (p *Storage) Ping() error {
	return p.conn.Ping(context.Background())
}
