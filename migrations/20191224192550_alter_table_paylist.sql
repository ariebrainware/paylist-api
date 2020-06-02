-- -- +goose Up
-- -- SQL in this section is executed when the migration is applied.
ALTER table paylists
ADD COLUMN due_date TIMESTAMP WITH TIME ZONE NULL;
-- -- +goose Down
ALTER table paylists
DROP COLUMN due_date;
-- -- SQL in this section is executed when the migration is rolled back.
