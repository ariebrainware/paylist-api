-- +goose Up
ALTER table paylists
ADD COLUMN due_date TIMESTAMP WITH TIME ZONE NULL;
-- +goose Down
ALTER table paylists
DROP COLUMN due_date;
