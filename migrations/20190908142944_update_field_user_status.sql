-- +goose Up
-- +goose StatementBegin
ALTER TABLE
  loggings CHANGE COLUMN userStatus user_status boolean;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
  loggings CHANGE COLUMN user_status userStatus boolean;
-- +goose StatementEnd