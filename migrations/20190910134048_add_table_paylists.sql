-- +goose Up
CREATE TABLE Paylists(
id serial NOT NULL,
name varchar NOT NULL,
amount int NOT NULL,
completed boolean,
username varchar NULL,
created_at timestamp with time zone,
deleted_at timestamp with time zone,
updated_at timestamp with time zone,
PRIMARY KEY (id)
);
-- +goose Down
DROP TABLE Paylists;