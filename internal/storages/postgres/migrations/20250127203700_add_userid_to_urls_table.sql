-- +goose Up
ALTER TABLE urls
    ADD COLUMN userid TEXT;

-- +goose Down
ALTER TABLE urls
    DROP COLUMN userid;
