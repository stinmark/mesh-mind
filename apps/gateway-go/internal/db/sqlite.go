package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./mesh_mind.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		tokens INTEGER DEFAULT 10
	);`

	if _, err := DB.Exec(query); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}
