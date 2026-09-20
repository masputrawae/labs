-- +goose Up
CREATE TABLE priorities (
	id    INTEGER PRIMARY KEY,
	name  TEXT NOT NULL UNIQUE,
	emoji TEXT NOT NULL
);

INSERT INTO priorities (id, name, emoji) VALUES
	(1, 'Highest', '🔴'),
	(2, 'High',    '🟠'),
	(3, 'Medium',  '🔵'),
	(4, 'Low',     '🟣'),
	(5, 'Lowest',  '⚫');


-- +goose Down
DROP TABLE priorities;
