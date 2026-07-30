//go:build integration

package identity_test

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/config"
	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/server"
	"github.com/fontis-dev/fontis-platform/runtime/identity/internal/store"
	pb "github.com/fontis-dev/fontis-platform/runtime/identity/proto"
)

func rpcCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func shortSockPath() string {
	b := make([]byte, 3)
	rand.Read(b)
	return filepath.Join(os.TempDir(), "fi_"+hex.EncodeToString(b)+".sock")
}

func newIdentityTestServer(t *testing.T) (pb.IdentityServiceClient, func()) {
	t.Helper()

	socketPath := shortSockPath()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	db.SetMaxOpenConns(1)

	st := store.New(db)
	cfg := &config.Config{SocketPath: socketPath}

	srv := server.New(cfg, st, nil)
	ctx := context.Background()

	if err := srv.Start(ctx); err != nil {
		t.Fatalf("start server: %v", err)
	}

	conn, err := grpc.Dial(
		socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", addr)
		}),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	client := pb.NewIdentityServiceClient(conn)

	return client, func() {
		conn.Close()
		srv.Shutdown(ctx)
		db.Close()
		os.Remove(socketPath)
	}
}

func TestIdentityIntegration_HouseholdCRUD(t *testing.T) {
	client, cleanup := newIdentityTestServer(t)
	defer cleanup()

	h, err := client.CreateHousehold(rpcCtx(t), &pb.CreateHouseholdRequest{Name: "Test Home"})
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}
	if h.Household.Name != "Test Home" {
		t.Errorf("got name %q, want %q", h.Household.Name, "Test Home")
	}
	if h.Household.Id == "" {
		t.Error("expected non-empty household id")
	}

	got, err := client.GetHousehold(rpcCtx(t), &pb.GetHouseholdRequest{HouseholdId: h.Household.Id})
	if err != nil {
		t.Fatalf("GetHousehold: %v", err)
	}
	if got.Household.Name != "Test Home" {
		t.Errorf("got name %q, want %q", got.Household.Name, "Test Home")
	}

	updated, err := client.UpdateHousehold(rpcCtx(t), &pb.UpdateHouseholdRequest{
		HouseholdId: h.Household.Id,
		Name:        "Updated Home",
	})
	if err != nil {
		t.Fatalf("UpdateHousehold: %v", err)
	}
	if updated.Household.Name != "Updated Home" {
		t.Errorf("got name %q, want %q", updated.Household.Name, "Updated Home")
	}

	list, err := client.ListHouseholds(rpcCtx(t), &pb.ListHouseholdsRequest{})
	if err != nil {
		t.Fatalf("ListHouseholds: %v", err)
	}
	if len(list.Households) != 1 {
		t.Errorf("got %d households, want 1", len(list.Households))
	}

	_, err = client.DeleteHousehold(rpcCtx(t), &pb.DeleteHouseholdRequest{HouseholdId: h.Household.Id})
	if err != nil {
		t.Fatalf("DeleteHousehold: %v", err)
	}

	_, err = client.GetHousehold(rpcCtx(t), &pb.GetHouseholdRequest{HouseholdId: h.Household.Id})
	if err == nil {
		t.Fatal("expected error after deleting household")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("got code %v, want %v", status.Code(err), codes.NotFound)
	}
}

func TestIdentityIntegration_ProfileCRUD(t *testing.T) {
	client, cleanup := newIdentityTestServer(t)
	defer cleanup()

	h, err := client.CreateHousehold(rpcCtx(t), &pb.CreateHouseholdRequest{Name: "Test"})
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	p, err := client.CreateProfile(rpcCtx(t), &pb.CreateProfileRequest{
		HouseholdId: h.Household.Id,
		DisplayName: "Alice",
		Role:        "admin",
	})
	if err != nil {
		t.Fatalf("CreateProfile: %v", err)
	}
	if p.Profile.DisplayName != "Alice" {
		t.Errorf("got display_name %q, want %q", p.Profile.DisplayName, "Alice")
	}
	if p.Profile.Role != "admin" {
		t.Errorf("got role %q, want %q", p.Profile.Role, "admin")
	}
	if p.Profile.HouseholdId != h.Household.Id {
		t.Errorf("got household_id %q, want %q", p.Profile.HouseholdId, h.Household.Id)
	}

	got, err := client.GetProfile(rpcCtx(t), &pb.GetProfileRequest{ProfileId: p.Profile.Id})
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if got.Profile.DisplayName != "Alice" {
		t.Errorf("got display_name %q, want %q", got.Profile.DisplayName, "Alice")
	}

	updated, err := client.UpdateProfile(rpcCtx(t), &pb.UpdateProfileRequest{
		ProfileId:   p.Profile.Id,
		DisplayName: "Alice Updated",
		Role:        "member",
	})
	if err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	if updated.Profile.DisplayName != "Alice Updated" {
		t.Errorf("got display_name %q, want %q", updated.Profile.DisplayName, "Alice Updated")
	}
	if updated.Profile.Role != "member" {
		t.Errorf("got role %q, want %q", updated.Profile.Role, "member")
	}

	list, err := client.ListProfiles(rpcCtx(t), &pb.ListProfilesRequest{HouseholdId: h.Household.Id})
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(list.Profiles) != 1 {
		t.Errorf("got %d profiles, want 1", len(list.Profiles))
	}

	_, err = client.DeleteProfile(rpcCtx(t), &pb.DeleteProfileRequest{ProfileId: p.Profile.Id})
	if err != nil {
		t.Fatalf("DeleteProfile: %v", err)
	}

	_, err = client.GetProfile(rpcCtx(t), &pb.GetProfileRequest{ProfileId: p.Profile.Id})
	if err == nil {
		t.Fatal("expected error after deleting profile")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("got code %v, want %v", status.Code(err), codes.NotFound)
	}
}

func TestIdentityIntegration_ProfileCascadeDelete(t *testing.T) {
	client, cleanup := newIdentityTestServer(t)
	defer cleanup()

	h, err := client.CreateHousehold(rpcCtx(t), &pb.CreateHouseholdRequest{Name: "Test"})
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	_, err = client.CreateProfile(rpcCtx(t), &pb.CreateProfileRequest{
		HouseholdId: h.Household.Id,
		DisplayName: "Alice",
	})
	if err != nil {
		t.Fatalf("CreateProfile: %v", err)
	}

	_, err = client.DeleteHousehold(rpcCtx(t), &pb.DeleteHouseholdRequest{HouseholdId: h.Household.Id})
	if err != nil {
		t.Fatalf("DeleteHousehold: %v", err)
	}

	list, err := client.ListProfiles(rpcCtx(t), &pb.ListProfilesRequest{HouseholdId: h.Household.Id})
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(list.Profiles) != 0 {
		t.Errorf("got %d profiles after cascade delete, want 0", len(list.Profiles))
	}
}

func TestIdentityIntegration_CreateProfileDefaultRole(t *testing.T) {
	client, cleanup := newIdentityTestServer(t)
	defer cleanup()

	h, err := client.CreateHousehold(rpcCtx(t), &pb.CreateHouseholdRequest{Name: "Test"})
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}

	p, err := client.CreateProfile(rpcCtx(t), &pb.CreateProfileRequest{
		HouseholdId: h.Household.Id,
		DisplayName: "Bob",
	})
	if err != nil {
		t.Fatalf("CreateProfile: %v", err)
	}
	if p.Profile.Role != "member" {
		t.Errorf("got role %q, want %q", p.Profile.Role, "member")
	}
}

func TestIdentityIntegration_Errors(t *testing.T) {
	client, cleanup := newIdentityTestServer(t)
	defer cleanup()

	_, err := client.GetHousehold(rpcCtx(t), &pb.GetHouseholdRequest{HouseholdId: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for nonexistent household")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("got code %v, want %v", status.Code(err), codes.NotFound)
	}

	_, err = client.CreateHousehold(rpcCtx(t), &pb.CreateHouseholdRequest{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("got code %v, want %v", status.Code(err), codes.InvalidArgument)
	}

	_, err = client.CreateProfile(rpcCtx(t), &pb.CreateProfileRequest{
		HouseholdId: "nonexistent",
		DisplayName: "Alice",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent household")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("got code %v, want %v", status.Code(err), codes.NotFound)
	}
}

func TestIdentityIntegration_ListEmpty(t *testing.T) {
	client, cleanup := newIdentityTestServer(t)
	defer cleanup()

	list, err := client.ListHouseholds(rpcCtx(t), &pb.ListHouseholdsRequest{})
	if err != nil {
		t.Fatalf("ListHouseholds: %v", err)
	}
	if len(list.Households) != 0 {
		t.Errorf("got %d households, want 0", len(list.Households))
	}
}
