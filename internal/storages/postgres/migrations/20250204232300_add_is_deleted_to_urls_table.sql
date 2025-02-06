-- +goose Up
ALTER TABLE urls
    ADD COLUMN deleted BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE urls
    DROP COLUMN deleted;
