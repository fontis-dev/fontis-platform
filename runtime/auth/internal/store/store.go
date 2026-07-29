package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Migrate(ctx context.Context) error {
	queries := []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS password_hashes (
			profile_id TEXT PRIMARY KEY,
			hash TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			profile_id TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			refresh_token TEXT NOT NULL UNIQUE,
			expires_at TEXT NOT NULL,
			refresh_expires_at TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS api_tokens (
			id TEXT PRIMARY KEY,
			profile_id TEXT NOT NULL,
			name TEXT NOT NULL,
			token_hash TEXT NOT NULL,
			created_at TEXT NOT NULL,
			expires_at TEXT
		)`,
	}

	for _, q := range queries {
		if _, err := s.db.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}

var (
	sessionTTL        = 15 * time.Minute
	refreshTTL        = 7 * 24 * time.Hour
)

const (
	argonTime    uint32 = 3
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 4
	argonKeyLen  uint32 = 32
	saltLen             = 16
	tokenLen            = 32
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, fmt.Errorf("generate random: %w", err)
	}
	return b, nil
}

func generateToken() (string, error) {
	b, err := generateRandomBytes(tokenLen)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashPassword(password string) (string, error) {
	salt, err := generateRandomBytes(saltLen)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonTime, argonThreads,
		hex.EncodeToString(salt), hex.EncodeToString(hash)), nil
}

func verifyPassword(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || !strings.HasPrefix(parts[3], "m=") {
		return false, fmt.Errorf("parse hash: invalid format")
	}

	var m, t uint32
	var p uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p)
	if err != nil {
		return false, fmt.Errorf("parse hash params: %w", err)
	}

	saltHex := parts[4]
	hashHex := parts[5]

	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false, fmt.Errorf("decode salt: %w", err)
	}
	expected, err := hex.DecodeString(hashHex)
	if err != nil {
		return false, fmt.Errorf("decode hash: %w", err)
	}

	actual := argon2.IDKey([]byte(password), salt, t, m, p, uint32(len(expected)))
	if len(actual) != len(expected) {
		return false, nil
	}

	for i := range actual {
		if actual[i] != expected[i] {
			return false, nil
		}
	}
	return true, nil
}



func newUUID() (string, error) {
	b, err := generateRandomBytes(16)
	if err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}
