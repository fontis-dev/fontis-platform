package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	pb "github.com/fontis-dev/fontis-platform/runtime/auth/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

type APIToken struct {
	ID        string
	ProfileID string
	Name      string
	TokenHash string
	CreatedAt time.Time
	ExpiresAt *time.Time
}

func (s *Store) CreateAPIToken(ctx context.Context, profileID, name string, expiresIn *time.Duration) (*APIToken, string, error) {
	id, err := newUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create api token: %w", err)
	}

	rawToken, err := generateToken()
	if err != nil {
		return nil, "", fmt.Errorf("create api token: %w", err)
	}

	tokenHash := hashToken(rawToken)
	now := time.Now().UTC()

	var expiresAt *time.Time
	if expiresIn != nil {
		t := now.Add(*expiresIn)
		expiresAt = &t
	}

	var expiresAtStr *string
	if expiresAt != nil {
		s := expiresAt.Format(time.RFC3339Nano)
		expiresAtStr = &s
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO api_tokens (id, profile_id, name, token_hash, created_at, expires_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		id, profileID, name, tokenHash,
		now.Format(time.RFC3339Nano), expiresAtStr)
	if err != nil {
		return nil, "", fmt.Errorf("create api token: %w", err)
	}

	return &APIToken{
		ID:        id,
		ProfileID: profileID,
		Name:      name,
		TokenHash: tokenHash,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}, rawToken, nil
}

func (s *Store) ValidateAPIToken(ctx context.Context, token string) (string, error) {
	tokenHash := hashToken(token)

	row := s.db.QueryRowContext(ctx,
		`SELECT profile_id, expires_at FROM api_tokens WHERE token_hash = ?`, tokenHash)

	var profileID string
	var expiresAtStr *string
	if err := row.Scan(&profileID, &expiresAtStr); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("validate token: not found")
		}
		return "", fmt.Errorf("validate token: %w", err)
	}

	if expiresAtStr != nil {
		expiresAt, err := time.Parse(time.RFC3339Nano, *expiresAtStr)
		if err != nil {
			return "", fmt.Errorf("validate token: %w", err)
		}
		if time.Now().UTC().After(expiresAt) {
			return "", fmt.Errorf("validate token: expired")
		}
	}

	return profileID, nil
}

func (s *Store) RevokeAPIToken(ctx context.Context, tokenID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM api_tokens WHERE id = ?`, tokenID)
	if err != nil {
		return fmt.Errorf("revoke api token: %w", err)
	}
	return nil
}

func (s *Store) ListAPITokens(ctx context.Context, profileID string) ([]*pb.ApiToken, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, profile_id, name, token_hash, created_at, expires_at
		 FROM api_tokens WHERE profile_id = ? ORDER BY created_at`, profileID)
	if err != nil {
		return nil, fmt.Errorf("list api tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*pb.ApiToken
	for rows.Next() {
		var id, pid, name, tokenHash, createdAt string
		var expiresAtStr *string
		if err := rows.Scan(&id, &pid, &name, &tokenHash, &createdAt, &expiresAtStr); err != nil {
			return nil, fmt.Errorf("list api tokens: %w", err)
		}

		ct, err := time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse created_at: %w", err)
		}

		t := &pb.ApiToken{
			Id:         id,
			ProfileId:  pid,
			Name:       name,
			TokenHash:  "",
			CreatedAt:  timestamppb.New(ct),
		}

		if expiresAtStr != nil {
			et, err := time.Parse(time.RFC3339Nano, *expiresAtStr)
			if err != nil {
				return nil, fmt.Errorf("parse expires_at: %w", err)
			}
			t.ExpiresAt = timestamppb.New(et)
		}

		tokens = append(tokens, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list api tokens: %w", err)
	}
	return tokens, nil
}

func apiTokenToProto(t *APIToken, rawToken string) *pb.ApiToken {
	pb := &pb.ApiToken{
		Id:         t.ID,
		ProfileId:  t.ProfileID,
		Name:       t.Name,
		CreatedAt:  timestamppb.New(t.CreatedAt),
	}
	if rawToken != "" {
		pb.TokenHash = rawToken
	}
	if t.ExpiresAt != nil {
		pb.ExpiresAt = timestamppb.New(*t.ExpiresAt)
	}
	return pb
}
