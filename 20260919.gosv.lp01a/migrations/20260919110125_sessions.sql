-- +goose Up
CREATE TABLE sessions (
	id         TEXT PRIMARY KEY,
	expires_at DATETIME NOT NULL,
	user_id    TEXT NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
);


-- +goose Down
DROP TABLE sessions;
