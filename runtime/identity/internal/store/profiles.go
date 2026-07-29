package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pb "github.com/fontis-dev/fontis-platform/runtime/identity/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateProfile(ctx context.Context, householdID, displayName, role string) (*pb.Profile, error) {
	id, err := newUUID()
	if err != nil {
		return nil, fmt.Errorf("create profile: %w", err)
	}
	now := time.Now().UTC()

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO profiles (id, household_id, display_name, avatar_url, role, created_at, updated_at)
		 VALUES (?, ?, ?, '', ?, ?, ?)`,
		id, householdID, displayName, role, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if err != nil {
		return nil, fmt.Errorf("create profile: %w", err)
	}

	return &pb.Profile{
		Id:          id,
		HouseholdId: householdID,
		DisplayName: displayName,
		Role:        role,
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}, nil
}

func (s *Store) GetProfile(ctx context.Context, profileID string) (*pb.Profile, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, household_id, display_name, avatar_url, role, created_at, updated_at
		 FROM profiles WHERE id = ?`, profileID)

	p, err := scanProfile(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("get profile: %w", err)
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return p, nil
}

func (s *Store) ListProfiles(ctx context.Context, householdID string) ([]*pb.Profile, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, household_id, display_name, avatar_url, role, created_at, updated_at
		 FROM profiles WHERE household_id = ? ORDER BY created_at`, householdID)
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}
	defer rows.Close()

	var profiles []*pb.Profile
	for rows.Next() {
		p, err := scanProfile(rows)
		if err != nil {
			return nil, fmt.Errorf("list profiles: %w", err)
		}
		profiles = append(profiles, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}
	return profiles, nil
}

func (s *Store) UpdateProfile(ctx context.Context, profileID, displayName, role string) (*pb.Profile, error) {
	now := time.Now().UTC()

	res, err := s.db.ExecContext(ctx,
		`UPDATE profiles SET display_name = ?, role = ?, updated_at = ? WHERE id = ?`,
		displayName, role, now.Format(time.RFC3339Nano), profileID)
	if err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}
	if n == 0 {
		return nil, fmt.Errorf("update profile: %w", sql.ErrNoRows)
	}

	return s.GetProfile(ctx, profileID)
}

func (s *Store) DeleteProfile(ctx context.Context, profileID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM profiles WHERE id = ?`, profileID)
	if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}
	return nil
}

func scanProfile(row scanner) (*pb.Profile, error) {
	var id, householdID, displayName, avatarURL, role, createdAt, updatedAt string
	if err := row.Scan(&id, &householdID, &displayName, &avatarURL, &role, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	ct, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse created_at: %w", err)
	}
	ut, err := time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse updated_at: %w", err)
	}

	return &pb.Profile{
		Id:          id,
		HouseholdId: householdID,
		DisplayName: displayName,
		AvatarUrl:   avatarURL,
		Role:        role,
		CreatedAt:   timestamppb.New(ct),
		UpdatedAt:   timestamppb.New(ut),
	}, nil
}
