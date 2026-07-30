package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Session struct {
	ID               string
	ProfileID        string
	Token            string
	RefreshToken     string
	ExpiresAt        time.Time
	RefreshExpiresAt time.Time
	CreatedAt        time.Time
}

func (s *Store) CreateSession(ctx context.Context, profileID string) (*Session, string, string, error) {
	id, err := newUUID()
	if err != nil {
		return nil, "", "", fmt.Errorf("create session: %w", err)
	}

	token, err := generateToken()
	if err != nil {
		return nil, "", "", fmt.Errorf("create session: %w", err)
	}

	refreshToken, err := generateToken()
	if err != nil {
		return nil, "", "", fmt.Errorf("create session: %w", err)
	}

	now := time.Now().UTC()
	expiresAt := now.Add(sessionTTL)
	refreshExpiresAt := now.Add(refreshTTL)

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO sessions (id, profile_id, token, refresh_token, expires_at, refresh_expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, profileID, token, refreshToken,
		expiresAt.Format(time.RFC3339Nano),
		refreshExpiresAt.Format(time.RFC3339Nano),
		now.Format(time.RFC3339Nano))
	if err != nil {
		return nil, "", "", fmt.Errorf("create session: %w", err)
	}

	return &Session{
		ID:               id,
		ProfileID:        profileID,
		Token:            token,
		RefreshToken:     refreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		CreatedAt:        now,
	}, token, refreshToken, nil
}

func (s *Store) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, profile_id, token, refresh_token, expires_at, refresh_expires_at, created_at
		 FROM sessions WHERE token = ?`, token)

	return scanSession(row)
}

func (s *Store) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*Session, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, profile_id, token, refresh_token, expires_at, refresh_expires_at, created_at
		 FROM sessions WHERE refresh_token = ?`, refreshToken)

	return scanSession(row)
}

func (s *Store) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE id = ?`, sessionID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}

func (s *Store) RefreshSession(ctx context.Context, oldRefreshToken string) (*Session, string, string, error) {
	sess, err := s.GetSessionByRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return nil, "", "", fmt.Errorf("refresh session: %w", err)
	}

	if time.Now().UTC().After(sess.RefreshExpiresAt) {
		_ = s.RevokeSession(ctx, sess.ID)
		return nil, "", "", fmt.Errorf("refresh session: expired refresh token")
	}

	if err := s.RevokeSession(ctx, sess.ID); err != nil {
		return nil, "", "", fmt.Errorf("refresh session: %w", err)
	}

	return s.CreateSession(ctx, sess.ProfileID)
}

func (s *Store) CleanExpiredSessions(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM sessions WHERE refresh_expires_at < ?`,
		time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("clean expired sessions: %w", err)
	}
	return nil
}

func scanSession(row scanner) (*Session, error) {
	var id, profileID, token, refreshToken, expiresAt, refreshExpiresAt, createdAt string
	if err := row.Scan(&id, &profileID, &token, &refreshToken, &expiresAt, &refreshExpiresAt, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("scan session: %w", err)
	}

	et, err := time.Parse(time.RFC3339Nano, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("parse expires_at: %w", err)
	}
	ret, err := time.Parse(time.RFC3339Nano, refreshExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("parse refresh_expires_at: %w", err)
	}
	ct, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse created_at: %w", err)
	}

	return &Session{
		ID:               id,
		ProfileID:        profileID,
		Token:            token,
		RefreshToken:     refreshToken,
		ExpiresAt:        et,
		RefreshExpiresAt: ret,
		CreatedAt:        ct,
	}, nil
}
