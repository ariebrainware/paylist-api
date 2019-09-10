-- +goose Up
CREATE TABLE Loggings(
 username varchar,
 token varchar,
 userStatus boolean,
 created_at timestamp with time zone,
 deleted_at timestamp with time zone
);
-- +goose Down
DROP TABLE Loggings;
