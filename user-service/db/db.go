package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "user.db")
	if err != nil {
		return nil, err
	}

	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL DEFAULT '',
		role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user'))
	);`

	_, err = database.Exec(query)
	return database, err
}