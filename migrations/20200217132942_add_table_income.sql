-- +goose Up
CREATE table Incomes (
id serial primary key not null,
username varchar not null,
balance int,
created_at timestamp with time zone,
deleted_at timestamp with time zone,
updated_at timestamp with time zone
);
-- +goose Down
DROP table Incomes;
