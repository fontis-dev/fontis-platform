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
	CREATE TABLE IF NOT EXISTS volumes (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		device_path TEXT NOT NULL,
		size_bytes INTEGER NOT NULL,
		encrypted INTEGER NOT NULL DEFAULT 0,
		filesystem TEXT NOT NULL DEFAULT 'btrfs',
		mount_point TEXT,
		state TEXT NOT NULL DEFAULT 'detached',
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS pools (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		volume_ids TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	`
	_, err := db.Exec(query)
	return err
}
