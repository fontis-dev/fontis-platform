package store

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	db.SetMaxOpenConns(1)

	st := New(db)
	if err := st.Migrate(context.Background()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return st
}

func TestCreateAndGetHousehold(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	h, err := st.CreateHousehold(ctx, "Test Home")
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}
	if h.Name != "Test Home" {
		t.Errorf("got name %q, want %q", h.Name, "Test Home")
	}
	if h.Id == "" {
		t.Error("expected non-empty household id")
	}

	got, err := st.GetHousehold(ctx, h.Id)
	if err != nil {
		t.Fatalf("GetHousehold: %v", err)
	}
	if got.Name != h.Name || got.Id != h.Id {
		t.Errorf("GetHousehold returned %+v, want %+v", got, h)
	}
}

func TestGetHouseholdNotFound(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, err := st.GetHousehold(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent household")
	}
}

func TestUpdateHousehold(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	h, err := st.CreateHousehold(ctx, "Original")
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	updated, err := st.UpdateHousehold(ctx, h.Id, "Updated")
	if err != nil {
		t.Fatalf("UpdateHousehold: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("got name %q, want %q", updated.Name, "Updated")
	}
}

func TestUpdateHouseholdNotFound(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, err := st.UpdateHousehold(ctx, "nonexistent", "Name")
	if err == nil {
		t.Fatal("expected error for nonexistent household")
	}
}

func TestDeleteHousehold(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	h, err := st.CreateHousehold(ctx, "To Delete")
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	if err := st.DeleteHousehold(ctx, h.Id); err != nil {
		t.Fatalf("DeleteHousehold: %v", err)
	}

	_, err = st.GetHousehold(ctx, h.Id)
	if err == nil {
		t.Error("expected error after deleting household")
	}
}

func TestListHouseholds(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.CreateHousehold(ctx, "A"); err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}
	if _, err := st.CreateHousehold(ctx, "B"); err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	households, err := st.ListHouseholds(ctx)
	if err != nil {
		t.Fatalf("ListHouseholds: %v", err)
	}
	if len(households) != 2 {
		t.Errorf("got %d households, want 2", len(households))
	}
}

func TestEmptyListHouseholds(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	households, err := st.ListHouseholds(ctx)
	if err != nil {
		t.Fatalf("ListHouseholds: %v", err)
	}
	if len(households) != 0 {
		t.Errorf("got %d households, want 0", len(households))
	}
}

func TestCreateAndGetProfile(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	h, err := st.CreateHousehold(ctx, "Test Home")
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	p, err := st.CreateProfile(ctx, h.Id, "Alice", "admin")
	if err != nil {
		t.Fatalf("CreateProfile: %v", err)
	}
	if p.DisplayName != "Alice" {
		t.Errorf("got display_name %q, want %q", p.DisplayName, "Alice")
	}
	if p.Role != "admin" {
		t.Errorf("got role %q, want %q", p.Role, "admin")
	}
	if p.HouseholdId != h.Id {
		t.Errorf("got household_id %q, want %q", p.HouseholdId, h.Id)
	}

	got, err := st.GetProfile(ctx, p.Id)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if got.DisplayName != p.DisplayName || got.Id != p.Id {
		t.Errorf("GetProfile returned %+v, want %+v", got, p)
	}
}

func TestGetProfileNotFound(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, err := st.GetProfile(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent profile")
	}
}

func TestListProfiles(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	h, err := st.CreateHousehold(ctx, "Test")
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	if _, err := st.CreateProfile(ctx, h.Id, "Alice", "admin"); err != nil {
		t.Fatalf("CreateProfile: %v", err)
	}
	if _, err := st.CreateProfile(ctx, h.Id, "Bob", "member"); err != nil {
		t.Fatalf("CreateProfile: %v", err)
	}

	profiles, err := st.ListProfiles(ctx, h.Id)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(profiles) != 2 {
		t.Errorf("got %d profiles, want 2", len(profiles))
	}
}

func TestListProfilesEmpty(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	h, err := st.CreateHousehold(ctx, "Test")
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	profiles, err := st.ListProfiles(ctx, h.Id)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("got %d profiles, want 0", len(profiles))
	}
}

func TestUpdateProfile(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	h, err := st.CreateHousehold(ctx, "Test")
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	p, err := st.CreateProfile(ctx, h.Id, "Alice", "member")
	if err != nil {
		t.Fatalf("CreateProfile: %v", err)
	}

	updated, err := st.UpdateProfile(ctx, p.Id, "Alice Updated", "admin")
	if err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	if updated.DisplayName != "Alice Updated" {
		t.Errorf("got display_name %q, want %q", updated.DisplayName, "Alice Updated")
	}
	if updated.Role != "admin" {
		t.Errorf("got role %q, want %q", updated.Role, "admin")
	}
}

func TestUpdateProfileNotFound(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, err := st.UpdateProfile(ctx, "nonexistent", "Name", "role")
	if err == nil {
		t.Fatal("expected error for nonexistent profile")
	}
}

func TestDeleteProfile(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	h, err := st.CreateHousehold(ctx, "Test")
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	p, err := st.CreateProfile(ctx, h.Id, "Alice", "member")
	if err != nil {
		t.Fatalf("CreateProfile: %v", err)
	}

	if err := st.DeleteProfile(ctx, p.Id); err != nil {
		t.Fatalf("DeleteProfile: %v", err)
	}

	_, err = st.GetProfile(ctx, p.Id)
	if err == nil {
		t.Error("expected error after deleting profile")
	}
}

func TestProfileCascadeDelete(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	h, err := st.CreateHousehold(ctx, "Test")
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	if _, err := st.CreateProfile(ctx, h.Id, "Alice", "member"); err != nil {
		t.Fatalf("CreateProfile: %v", err)
	}

	if err := st.DeleteHousehold(ctx, h.Id); err != nil {
		t.Fatalf("DeleteHousehold: %v", err)
	}

	profiles, err := st.ListProfiles(ctx, h.Id)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("got %d profiles after household delete, want 0", len(profiles))
	}
}

func TestNewUUID(t *testing.T) {
	u1, err := newUUID()
	if err != nil {
		t.Fatalf("newUUID: %v", err)
	}
	u2, err := newUUID()
	if err != nil {
		t.Fatalf("newUUID: %v", err)
	}
	if u1 == u2 {
		t.Error("expected different UUIDs")
	}
	if len(u1) != 36 {
		t.Errorf("expected UUID length 36, got %d", len(u1))
	}
}
