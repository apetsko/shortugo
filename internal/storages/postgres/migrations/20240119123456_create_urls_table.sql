-- +goose Up
CREATE TABLE IF NOT EXISTS urls (
    id TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    date DATE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS urls;
