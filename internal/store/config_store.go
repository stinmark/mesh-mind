package store

import (
	"database/sql"
)

type ConfigStore struct {
	db *sql.DB
}

func NewConfigStore(db *sql.DB) *ConfigStore {
	return &ConfigStore{db: db}
}

func (s *ConfigStore) Set(key, value string) error {
	query := `INSERT INTO settings (key, value) VALUES (?, ?) 
	          ON CONFLICT(key) DO UPDATE SET value=excluded.value;`
	_, err := s.db.Exec(query, key, value)
	return err
}

func (s *ConfigStore) Get(key string) (string, error) {
	var value string
	query := `SELECT value FROM settings WHERE key = ?;`
	err := s.db.QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}
