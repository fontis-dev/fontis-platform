//go:build integration

package auth_test

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

	"github.com/fontis-dev/fontis-platform/runtime/auth/internal/config"
	"github.com/fontis-dev/fontis-platform/runtime/auth/internal/server"
	"github.com/fontis-dev/fontis-platform/runtime/auth/internal/store"
	pb "github.com/fontis-dev/fontis-platform/runtime/auth/proto"
)

func shortSockPath() string {
	b := make([]byte, 3)
	rand.Read(b)
	return filepath.Join(os.TempDir(), "fa_"+hex.EncodeToString(b)+".sock")
}

func newAuthTestServer(t *testing.T) (pb.AuthServiceClient, *store.Store, func()) {
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

	client := pb.NewAuthServiceClient(conn)

	return client, st, func() {
		conn.Close()
		srv.Shutdown(ctx)
		st.Close()
		os.Remove(socketPath)
	}
}

func TestAuthIntegration_SessionLifecycle(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sess, err := client.CreateSession(ctx, &pb.CreateSessionRequest{ProfileId: "profile-1"})
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if sess.Session.Id == "" {
		t.Error("expected non-empty session id")
	}
	if sess.Session.Token == "" {
		t.Error("expected non-empty token")
	}
	if sess.Session.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}

	valid, err := client.ValidateSession(ctx, &pb.ValidateSessionRequest{Token: sess.Session.Token})
	if err != nil {
		t.Fatalf("ValidateSession: %v", err)
	}
	if !valid.Valid {
		t.Error("expected session to be valid")
	}
	if valid.ProfileId != "profile-1" {
		t.Errorf("got profile_id %q, want %q", valid.ProfileId, "profile-1")
	}

	_, err = client.RevokeSession(ctx, &pb.RevokeSessionRequest{SessionId: sess.Session.Id})
	if err != nil {
		t.Fatalf("RevokeSession: %v", err)
	}

	valid, err = client.ValidateSession(ctx, &pb.ValidateSessionRequest{Token: sess.Session.Token})
	if err != nil {
		t.Fatalf("ValidateSession after revoke: %v", err)
	}
	if valid.Valid {
		t.Error("expected session to be invalid after revoke")
	}
}

func TestAuthIntegration_SessionRefresh(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sess, err := client.CreateSession(ctx, &pb.CreateSessionRequest{ProfileId: "profile-1"})
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	refreshed, err := client.RefreshSession(ctx, &pb.RefreshSessionRequest{
		RefreshToken: sess.Session.RefreshToken,
	})
	if err != nil {
		t.Fatalf("RefreshSession: %v", err)
	}
	if refreshed.Session.Token == sess.Session.Token {
		t.Error("expected new token different from old")
	}
	if refreshed.Session.RefreshToken == sess.Session.RefreshToken {
		t.Error("expected new refresh token different from old")
	}

	valid, err := client.ValidateSession(ctx, &pb.ValidateSessionRequest{Token: refreshed.Session.Token})
	if err != nil {
		t.Fatalf("ValidateSession new token: %v", err)
	}
	if !valid.Valid {
		t.Error("expected new session to be valid")
	}

	valid, err = client.ValidateSession(ctx, &pb.ValidateSessionRequest{Token: sess.Session.Token})
	if err != nil {
		t.Fatalf("ValidateSession old token: %v", err)
	}
	if valid.Valid {
		t.Error("expected old session token to be invalid after refresh")
	}
}

func TestAuthIntegration_Authenticate(t *testing.T) {
	client, st, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := st.SetPassword(ctx, "profile-1", "correct-password"); err != nil {
		t.Fatalf("SetPassword: %v", err)
	}

	sess, err := client.Authenticate(ctx, &pb.AuthenticateRequest{
		ProfileId: "profile-1",
		Password:  "correct-password",
	})
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if sess.Session.ProfileId != "profile-1" {
		t.Errorf("got profile_id %q, want %q", sess.Session.ProfileId, "profile-1")
	}
	if sess.Session.Token == "" {
		t.Error("expected non-empty session token")
	}
}

func TestAuthIntegration_AuthenticateInvalidPassword(t *testing.T) {
	client, st, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := st.SetPassword(ctx, "profile-1", "correct-password"); err != nil {
		t.Fatalf("SetPassword: %v", err)
	}

	_, err := client.Authenticate(ctx, &pb.AuthenticateRequest{
		ProfileId: "profile-1",
		Password:  "wrong-password",
	})
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("got code %v, want %v", status.Code(err), codes.PermissionDenied)
	}
}

func TestAuthIntegration_APITokenLifecycle(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tok, err := client.CreateApiToken(ctx, &pb.CreateApiTokenRequest{
		ProfileId: "profile-1",
		Name:      "test-token",
	})
	if err != nil {
		t.Fatalf("CreateApiToken: %v", err)
	}
	if tok.Token.TokenHash == "" {
		t.Error("expected token_hash (raw token) in create response")
	}
	if tok.Token.Id == "" {
		t.Error("expected non-empty token id")
	}

	valid, err := client.ValidateApiToken(ctx, &pb.ValidateApiTokenRequest{
		Token: tok.Token.TokenHash,
	})
	if err != nil {
		t.Fatalf("ValidateApiToken: %v", err)
	}
	if !valid.Valid {
		t.Error("expected API token to be valid")
	}
	if valid.ProfileId != "profile-1" {
		t.Errorf("got profile_id %q, want %q", valid.ProfileId, "profile-1")
	}

	_, err = client.RevokeApiToken(ctx, &pb.RevokeApiTokenRequest{TokenId: tok.Token.Id})
	if err != nil {
		t.Fatalf("RevokeApiToken: %v", err)
	}

	valid, err = client.ValidateApiToken(ctx, &pb.ValidateApiTokenRequest{
		Token: tok.Token.TokenHash,
	})
	if err != nil {
		t.Fatalf("ValidateApiToken after revoke: %v", err)
	}
	if valid.Valid {
		t.Error("expected API token to be invalid after revoke")
	}
}

func TestAuthIntegration_ListAPITokens(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.CreateApiToken(ctx, &pb.CreateApiTokenRequest{
		ProfileId: "profile-1",
		Name:      "token-a",
	})
	if err != nil {
		t.Fatalf("CreateApiToken: %v", err)
	}

	_, err = client.CreateApiToken(ctx, &pb.CreateApiTokenRequest{
		ProfileId: "profile-1",
		Name:      "token-b",
	})
	if err != nil {
		t.Fatalf("CreateApiToken: %v", err)
	}

	list, err := client.ListApiTokens(ctx, &pb.ListApiTokensRequest{ProfileId: "profile-1"})
	if err != nil {
		t.Fatalf("ListApiTokens: %v", err)
	}
	if len(list.Tokens) != 2 {
		t.Errorf("got %d tokens, want 2", len(list.Tokens))
	}
	for _, tok := range list.Tokens {
		if tok.TokenHash != "" {
			t.Error("expected empty token_hash in list response")
		}
	}
}

func TestAuthIntegration_ValidateInvalidSession(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	valid, err := client.ValidateSession(ctx, &pb.ValidateSessionRequest{
		Token: "nonexistent-token",
	})
	if err != nil {
		t.Fatalf("ValidateSession: %v", err)
	}
	if valid.Valid {
		t.Error("expected invalid session for nonexistent token")
	}
}

func TestAuthIntegration_ValidateInvalidAPIToken(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	valid, err := client.ValidateApiToken(ctx, &pb.ValidateApiTokenRequest{
		Token: "invalid-token",
	})
	if err != nil {
		t.Fatalf("ValidateApiToken: %v", err)
	}
	if valid.Valid {
		t.Error("expected invalid API token for nonexistent token")
	}
}

func TestAuthIntegration_MissingFields(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.CreateSession(ctx, &pb.CreateSessionRequest{ProfileId: ""})
	if err == nil {
		t.Fatal("expected error for empty profile_id")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("got code %v, want %v", status.Code(err), codes.InvalidArgument)
	}

	_, err = client.Authenticate(ctx, &pb.AuthenticateRequest{ProfileId: "", Password: ""})
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("got code %v, want %v", status.Code(err), codes.InvalidArgument)
	}

	_, err = client.CreateApiToken(ctx, &pb.CreateApiTokenRequest{ProfileId: "", Name: ""})
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("got code %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}

func TestAuthIntegration_RefreshInvalidToken(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.RefreshSession(ctx, &pb.RefreshSessionRequest{RefreshToken: "invalid"})
	if err == nil {
		t.Fatal("expected error for invalid refresh token")
	}
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("got code %v, want %v", status.Code(err), codes.Unauthenticated)
	}
}

func TestAuthIntegration_AuthenticateNonexistentProfile(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.Authenticate(ctx, &pb.AuthenticateRequest{
		ProfileId: "nonexistent",
		Password:  "password",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent profile")
	}
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("got code %v, want %v", status.Code(err), codes.PermissionDenied)
	}
}

func TestAuthIntegration_RevokeNonexistent(t *testing.T) {
	client, _, cleanup := newAuthTestServer(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.RevokeSession(ctx, &pb.RevokeSessionRequest{SessionId: "nonexistent"})
	if err != nil {
		t.Fatalf("RevokeSession nonexistent: %v", err)
	}

	_, err = client.RevokeApiToken(ctx, &pb.RevokeApiTokenRequest{TokenId: "nonexistent"})
	if err != nil {
		t.Fatalf("RevokeApiToken nonexistent: %v", err)
	}
}
