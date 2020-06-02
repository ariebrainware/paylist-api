-- -- +goose Up
-- -- +goose StatementBegin
ALTER TABLE Users
ALTER COLUMN created_at TYPE timestamp without time zone;
-- -- +goose StatementEnd
-- -- +goose Down
-- -- +goose StatementBegin
ALTER TABLE Users 
ALTER COLUMN created_at TYPE timestamp with time zone;
-- -- +goose StatementEnd
