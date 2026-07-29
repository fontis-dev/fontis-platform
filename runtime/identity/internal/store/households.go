package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pb "github.com/fontis-dev/fontis-platform/runtime/identity/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateHousehold(ctx context.Context, name string) (*pb.Household, error) {
	id, err := newUUID()
	if err != nil {
		return nil, fmt.Errorf("create household: %w", err)
	}
	now := time.Now().UTC()

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO households (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		id, name, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if err != nil {
		return nil, fmt.Errorf("create household: %w", err)
	}

	return &pb.Household{
		Id:        id,
		Name:      name,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}, nil
}

func (s *Store) GetHousehold(ctx context.Context, householdID string) (*pb.Household, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, created_at, updated_at FROM households WHERE id = ?`, householdID)

	h, err := scanHousehold(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("get household: %w", err)
		}
		return nil, fmt.Errorf("get household: %w", err)
	}
	return h, nil
}

func (s *Store) UpdateHousehold(ctx context.Context, householdID, name string) (*pb.Household, error) {
	now := time.Now().UTC()

	res, err := s.db.ExecContext(ctx,
		`UPDATE households SET name = ?, updated_at = ? WHERE id = ?`,
		name, now.Format(time.RFC3339Nano), householdID)
	if err != nil {
		return nil, fmt.Errorf("update household: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update household: %w", err)
	}
	if n == 0 {
		return nil, fmt.Errorf("update household: %w", sql.ErrNoRows)
	}

	return s.GetHousehold(ctx, householdID)
}

func (s *Store) DeleteHousehold(ctx context.Context, householdID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM households WHERE id = ?`, householdID)
	if err != nil {
		return fmt.Errorf("delete household: %w", err)
	}
	return nil
}

func (s *Store) ListHouseholds(ctx context.Context) ([]*pb.Household, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, created_at, updated_at FROM households ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list households: %w", err)
	}
	defer rows.Close()

	var households []*pb.Household
	for rows.Next() {
		h, err := scanHousehold(rows)
		if err != nil {
			return nil, fmt.Errorf("list households: %w", err)
		}
		households = append(households, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list households: %w", err)
	}
	return households, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanHousehold(row scanner) (*pb.Household, error) {
	var id, name, createdAt, updatedAt string
	if err := row.Scan(&id, &name, &createdAt, &updatedAt); err != nil {
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

	return &pb.Household{
		Id:        id,
		Name:      name,
		CreatedAt: timestamppb.New(ct),
		UpdatedAt: timestamppb.New(ut),
	}, nil
}
