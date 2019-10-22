-- +goose Up
-- +goose StatementBegin
ALTER TABLE Paylists ADD FOREIGN KEY(username) REFERENCES Users(username);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE Paylists DROP FOREIGN KEY(username);
-- +goose StatementEnd
