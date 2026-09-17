package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitSQLite3(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
