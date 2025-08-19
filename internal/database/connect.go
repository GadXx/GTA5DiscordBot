package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Connect(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS inactive (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT,
		user_name TEXT,
		reason TEXT,
		created_at TIMESTAMP,
		end_at TIMESTAMP
	);
	`
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
