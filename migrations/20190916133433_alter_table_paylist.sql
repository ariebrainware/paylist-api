-- +goose Up
ALTER TABLE Paylists ADD FOREIGN KEY(username) REFERENCES Users(username);
-- +goose Down
ALTER TABLE Paylists DROP FOREIGN KEY(username);
