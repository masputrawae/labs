package database

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitSQLite3(ctx context.Context, file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	q := "PRAGMA foreign_keys=ON"
	if _, err := db.ExecContext(ctx, q); err != nil {
		return nil, err
	}

	return db, nil
}
