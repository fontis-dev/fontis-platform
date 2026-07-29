package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

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

func TestSetAndVerifyPassword(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if err := st.SetPassword(ctx, "profile-1", "correct-password"); err != nil {
		t.Fatalf("SetPassword: %v", err)
	}

	ok, err := st.VerifyPassword(ctx, "profile-1", "correct-password")
	if err != nil {
		t.Fatalf("VerifyPassword: %v", err)
	}
	if !ok {
		t.Error("expected password verification to succeed")
	}

	ok, err = st.VerifyPassword(ctx, "profile-1", "wrong-password")
	if err != nil {
		t.Fatalf("VerifyPassword: %v", err)
	}
	if ok {
		t.Error("expected password verification to fail")
	}
}

func TestSetPasswordUpdatesExisting(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if err := st.SetPassword(ctx, "profile-1", "first"); err != nil {
		t.Fatalf("SetPassword: %v", err)
	}
	if err := st.SetPassword(ctx, "profile-1", "second"); err != nil {
		t.Fatalf("SetPassword: %v", err)
	}

	ok, err := st.VerifyPassword(ctx, "profile-1", "second")
	if err != nil {
		t.Fatalf("VerifyPassword: %v", err)
	}
	if !ok {
		t.Error("expected new password to work")
	}

	ok, err = st.VerifyPassword(ctx, "profile-1", "first")
	if err != nil {
		t.Fatalf("VerifyPassword: %v", err)
	}
	if ok {
		t.Error("expected old password to fail after update")
	}
}

func TestVerifyPasswordNotFound(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, err := st.VerifyPassword(ctx, "nonexistent", "password")
	if err == nil {
		t.Fatal("expected error for nonexistent profile")
	}
}

func TestHasPassword(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	exists, err := st.HasPassword(ctx, "profile-1")
	if err != nil {
		t.Fatalf("HasPassword: %v", err)
	}
	if exists {
		t.Error("expected no password initially")
	}

	if err := st.SetPassword(ctx, "profile-1", "pass"); err != nil {
		t.Fatalf("SetPassword: %v", err)
	}

	exists, err = st.HasPassword(ctx, "profile-1")
	if err != nil {
		t.Fatalf("HasPassword: %v", err)
	}
	if !exists {
		t.Error("expected password to exist after set")
	}
}

func TestCreateAndValidateSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	sess, token, refreshToken, err := st.CreateSession(ctx, "profile-1")
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if sess.ID == "" {
		t.Error("expected non-empty session id")
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	if refreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if sess.ProfileID != "profile-1" {
		t.Errorf("got profile_id %q, want %q", sess.ProfileID, "profile-1")
	}

	got, err := st.GetSessionByToken(ctx, token)
	if err != nil {
		t.Fatalf("GetSessionByToken: %v", err)
	}
	if got.ID != sess.ID {
		t.Errorf("got session id %q, want %q", got.ID, sess.ID)
	}
}

func TestValidateExpiredSession(t *testing.T) {
	oldTTL := sessionTTL
	sessionTTL = -1 * time.Second
	defer func() { sessionTTL = oldTTL }()

	st := newTestStore(t)
	ctx := context.Background()

	_, token, _, err := st.CreateSession(ctx, "profile-1")
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	got, err := st.GetSessionByToken(ctx, token)
	if err != nil {
		t.Fatalf("GetSessionByToken: %v", err)
	}

	if time.Now().UTC().Before(got.ExpiresAt) {
		t.Error("expected session to be expired")
	}
}

func TestRevokeSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, token, _, err := st.CreateSession(ctx, "profile-1")
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	sess, err := st.GetSessionByToken(ctx, token)
	if err != nil {
		t.Fatalf("GetSessionByToken: %v", err)
	}

	if err := st.RevokeSession(ctx, sess.ID); err != nil {
		t.Fatalf("RevokeSession: %v", err)
	}

	_, err = st.GetSessionByToken(ctx, token)
	if err == nil {
		t.Error("expected error after revoking session")
	}
}

func TestRefreshSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, _, refreshToken, err := st.CreateSession(ctx, "profile-1")
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	newSess, newToken, newRefreshToken, err := st.RefreshSession(ctx, refreshToken)
	if err != nil {
		t.Fatalf("RefreshSession: %v", err)
	}
	if newSess.ProfileID != "profile-1" {
		t.Errorf("got profile_id %q, want %q", newSess.ProfileID, "profile-1")
	}
	if newToken == "" {
		t.Error("expected new token")
	}
	if newRefreshToken == "" || newRefreshToken == refreshToken {
		t.Error("expected new refresh token, different from old")
	}

	oldSess, err := st.GetSessionByRefreshToken(ctx, refreshToken)
	if err == nil {
		t.Errorf("expected old refresh token to be revoked, got session %s", oldSess.ID)
	}
}

func TestRefreshExpiredToken(t *testing.T) {
	oldRefreshTTL := refreshTTL
	refreshTTL = -1 * time.Second
	defer func() { refreshTTL = oldRefreshTTL }()

	st := newTestStore(t)
	ctx := context.Background()

	_, _, refreshToken, err := st.CreateSession(ctx, "profile-1")
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	_, _, _, err = st.RefreshSession(ctx, refreshToken)
	if err == nil {
		t.Error("expected error for expired refresh token")
	}
}

func TestCreateAndValidateAPIToken(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, rawToken, err := st.CreateAPIToken(ctx, "profile-1", "test-token", nil)
	if err != nil {
		t.Fatalf("CreateAPIToken: %v", err)
	}
	if rawToken == "" {
		t.Error("expected non-empty raw token")
	}

	profileID, err := st.ValidateAPIToken(ctx, rawToken)
	if err != nil {
		t.Fatalf("ValidateAPIToken: %v", err)
	}
	if profileID != "profile-1" {
		t.Errorf("got profile_id %q, want %q", profileID, "profile-1")
	}
}

func TestValidateExpiredAPIToken(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	expiresIn := -1 * time.Hour
	_, rawToken, err := st.CreateAPIToken(ctx, "profile-1", "expired", &expiresIn)
	if err != nil {
		t.Fatalf("CreateAPIToken: %v", err)
	}

	_, err = st.ValidateAPIToken(ctx, rawToken)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestValidateInvalidAPIToken(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, err := st.ValidateAPIToken(ctx, "invalid-token")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestRevokeAPIToken(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	tok, rawToken, err := st.CreateAPIToken(ctx, "profile-1", "revocable", nil)
	if err != nil {
		t.Fatalf("CreateAPIToken: %v", err)
	}

	if err := st.RevokeAPIToken(ctx, tok.ID); err != nil {
		t.Fatalf("RevokeAPIToken: %v", err)
	}

	_, err = st.ValidateAPIToken(ctx, rawToken)
	if err == nil {
		t.Error("expected error after revoking token")
	}
}

func TestListAPITokens(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, _, err := st.CreateAPIToken(ctx, "profile-1", "token-1", nil); err != nil {
		t.Fatalf("CreateAPIToken: %v", err)
	}
	if _, _, err := st.CreateAPIToken(ctx, "profile-1", "token-2", nil); err != nil {
		t.Fatalf("CreateAPIToken: %v", err)
	}
	if _, _, err := st.CreateAPIToken(ctx, "other-profile", "token-3", nil); err != nil {
		t.Fatalf("CreateAPIToken: %v", err)
	}

	tokens, err := st.ListAPITokens(ctx, "profile-1")
	if err != nil {
		t.Fatalf("ListAPITokens: %v", err)
	}
	if len(tokens) != 2 {
		t.Errorf("got %d tokens, want 2", len(tokens))
	}

	for _, tok := range tokens {
		if tok.TokenHash != "" {
			t.Error("expected empty token_hash in list response")
		}
	}
}

func TestCleanExpiredSessions(t *testing.T) {
	oldRefreshTTL := refreshTTL
	refreshTTL = -1 * time.Hour
	defer func() { refreshTTL = oldRefreshTTL }()

	st := newTestStore(t)
	ctx := context.Background()

	_, token, _, err := st.CreateSession(ctx, "profile-1")
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	if err := st.CleanExpiredSessions(ctx); err != nil {
		t.Fatalf("CleanExpiredSessions: %v", err)
	}

	_, err = st.GetSessionByToken(ctx, token)
	if err == nil {
		t.Error("expected session to be cleaned")
	}
}

func TestHashTokenDeterministic(t *testing.T) {
	h1 := hashToken("test-token")
	h2 := hashToken("test-token")
	if h1 != h2 {
		t.Error("expected hash to be deterministic")
	}

	h3 := hashToken("different")
	if h1 == h3 {
		t.Error("expected different tokens to produce different hashes")
	}
}

func TestGenerateTokenUnique(t *testing.T) {
	t1, err := generateToken()
	if err != nil {
		t.Fatalf("generateToken: %v", err)
	}
	t2, err := generateToken()
	if err != nil {
		t.Fatalf("generateToken: %v", err)
	}
	if t1 == t2 {
		t.Error("expected unique tokens")
	}
}
