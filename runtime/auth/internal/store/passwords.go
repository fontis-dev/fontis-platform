package store

import (
	"context"
	"fmt"
	"time"
)

func (s *Store) SetPassword(ctx context.Context, profileID, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("set password: %w", err)
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO password_hashes (profile_id, hash, created_at)
		 VALUES (?, ?, ?)`,
		profileID, hash, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("set password: %w", err)
	}

	return nil
}

func (s *Store) VerifyPassword(ctx context.Context, profileID, password string) (bool, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT hash FROM password_hashes WHERE profile_id = ?`, profileID)

	var encoded string
	if err := row.Scan(&encoded); err != nil {
		return false, fmt.Errorf("verify password: %w", err)
	}

	ok, err := verifyPassword(password, encoded)
	if err != nil {
		return false, fmt.Errorf("verify password: %w", err)
	}
	return ok, nil
}

func (s *Store) HasPassword(ctx context.Context, profileID string) (bool, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM password_hashes WHERE profile_id = ?`, profileID)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("has password: %w", err)
	}
	return count > 0, nil
}
