-- +goose Up
ALTER table incomes
RENAME COLUMN balance TO income;
-- +goose Down
ALTER table incomes
RENAME COLUMN income TO balance;
