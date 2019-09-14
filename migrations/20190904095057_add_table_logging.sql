-- +goose Up
CREATE TABLE loggings (
    username varchar(255),
    token varchar(255),
    userStatus boolean,
    created_at timestamp null,
    deleted_at timestamp null
);
-- +goose Down
DROP TABLE loggings;
