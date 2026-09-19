-- +goose Up
CREATE TABLE statuses (
	id    INTEGER PRIMARY KEY,
	name  TEXT NOT NULL UNIQUE,
	emoji TEXT NOT NULL
);

INSERT INTO statuses (id, name, emoji) VALUES
	(1, 'Planning',    '📝'),
	(2, 'In Progress', '⏳'),
	(3, 'Paused',      '⏸️'),
	(4, 'Done',        '✅'),
	(5, 'Cancelled',   '❌');


-- +goose Down
DROP TABLE statuses;
