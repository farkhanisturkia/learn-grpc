package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "job.db")
	if err != nil {
		return nil, err
	}

	query := `CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		company TEXT NOT NULL
	);`

	_, err = database.Exec(query)
	return database, err
}