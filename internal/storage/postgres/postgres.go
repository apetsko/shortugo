package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/apetsko/shortugo/internal/storage/shared"
	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	conn *pgx.Conn
}

func New(connString string) (*Postgres, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
		shortUrl TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		date DATE NOT NULL
	);`

	_, err = conn.Exec(context.Background(), createTableQuery)
	if err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Postgres{conn: conn}, nil
}

func (p *Postgres) Close() error {
	return p.conn.Close(context.Background())
}

func (p *Postgres) Put(shortURL, url string) error {
	currentDate := time.Now().Format(time.RFC3339)
	log.Println(shortURL, url, currentDate)
	insertQuery := `
	INSERT INTO urls (shortUrl, url, date)
	VALUES ($1, $2, $3)
	ON CONFLICT (shortUrl)
	DO UPDATE SET date = EXCLUDED.date;`

	_, err := p.conn.Exec(context.Background(), insertQuery, shortURL, url, currentDate)
	if err != nil {
		return fmt.Errorf("failed to insert URL: %w", err)
	}

	return nil
}

func (p *Postgres) Get(shortURL string) (string, error) {
	var url string
	err := p.conn.QueryRow(context.Background(), "SELECT url FROM urls WHERE shortUrl=$1", shortURL).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("URL not found: %s. %w", shortURL, shared.ErrNotFound)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	return url, nil
}

func (p *Postgres) Ping() error {
	return p.conn.Ping(context.Background())
}
