-- +goose Up
SELECT 'CREATE DATABASE practicum'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'practicum')
LIMIT 1;

-- +goose Down
DROP DATABASE IF EXISTS practicum;
