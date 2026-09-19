-- +goose Up
CREATE TABLE users (
	id         TEXT PRIMARY KEY,
	name       TEXT NOT NULL,
	email      TEXT NOT NULL UNIQUE,
	username   TEXT NOT NULL UNIQUE,
	password   TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);


-- +goose Down
DROP TABLE users;
