package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "education.db")
	if err != nil {
		return nil, err
	}

	query := `CREATE TABLE IF NOT EXISTS educations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		level TEXT NOT NULL,
		program TEXT NOT NULL,
		university TEXT NOT NULL
	);`

	_, err = database.Exec(query)
	return database, err
}