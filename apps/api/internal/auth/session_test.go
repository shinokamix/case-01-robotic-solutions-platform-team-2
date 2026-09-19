package auth

import (
	"testing"
	"time"
)

func TestRefreshSessionTimeouts(t *testing.T) {
	now := time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)
	session, token, err := NewSession("user-1", now)
	if err != nil {
		t.Fatal(err)
	}
	if len(token) == 0 || len(session.TokenHash) != 32 {
		t.Fatal("session token was not generated")
	}

	refreshed, valid := RefreshSession(session, now.Add(13*24*time.Hour))
	if !valid {
		t.Fatal("active session must be valid")
	}
	if !refreshed.IdleExpiresAt.Equal(now.Add(27 * 24 * time.Hour)) {
		t.Fatalf("unexpected refreshed idle expiry: %s", refreshed.IdleExpiresAt)
	}
	refreshed, valid = RefreshSession(refreshed, now.Add(20*24*time.Hour))
	if !valid || !refreshed.IdleExpiresAt.Equal(session.AbsoluteExpiresAt) {
		t.Fatalf("idle expiry must be capped by absolute expiry: %s", refreshed.IdleExpiresAt)
	}
	if _, valid := RefreshSession(session, now.Add(14*24*time.Hour)); valid {
		t.Fatal("session must expire after 14 days of inactivity")
	}

	session.IdleExpiresAt = now.Add(40 * 24 * time.Hour)
	if _, valid := RefreshSession(session, now.Add(30*24*time.Hour)); valid {
		t.Fatal("session must expire after absolute 30 day lifetime")
	}
}
