-- +goose Up
CREATE TABLE Users(
id serial NOT NULL,
email varchar NOT NULL,
name varchar NOT NULL,
username varchar NOT NULL,
password varchar NOT NULL,
balance int NULL,
created_at timestamp with time zone,
deleted_at timestamp with time zone,
updated_at timestamp with time zone,
PRIMARY KEY(id, username)
);
-- +goose Down
DROP TABLE Users;

