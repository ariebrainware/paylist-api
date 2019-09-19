-- +goose Up
CREATE TABLE Paylists(
id serial NOT NULL PRIMARY KEY,
name varchar NOT NULL,
amount int NOT NULL,
completed boolean,
username varchar NULL,
created_at timestamp with time zone,
deleted_at timestamp with time zone,
updated_at timestamp with time zone
);
-- +goose Down
DROP TABLE Paylists;