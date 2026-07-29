package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS networks (
		id TEXT PRIMARY KEY,
		ssid TEXT NOT NULL,
		security_type TEXT NOT NULL DEFAULT 'wpa2',
		priority INTEGER NOT NULL DEFAULT 0,
		auto_connect INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS interfaces (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		type TEXT NOT NULL,
		mac_address TEXT,
		state TEXT NOT NULL DEFAULT 'down'
	);
	`
	_, err := db.Exec(query)
	return err
}
