-- +goose Up
ALTER table loggings
RENAME COLUMN userStatus TO user_status;
-- +goose Down
ALTER table loggings
RENAME COLUMN user_status TO userStatus;
