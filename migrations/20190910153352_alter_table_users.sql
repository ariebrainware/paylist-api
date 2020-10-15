-- +goose Up
ALTER TABLE Users
ALTER COLUMN created_at TYPE timestamp without time zone;
-- +goose Down
ALTER TABLE Users 
ALTER COLUMN created_at TYPE timestamp with time zone;
