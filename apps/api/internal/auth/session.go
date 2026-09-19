package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	idleSessionLifetime     = 14 * 24 * time.Hour
	absoluteSessionLifetime = 30 * 24 * time.Hour
	sessionTokenBytes       = 32
)

func NewSession(userID string, now time.Time) (Session, string, error) {
	raw := make([]byte, sessionTokenBytes)
	if _, err := rand.Read(raw); err != nil {
		return Session{}, "", fmt.Errorf("generate session token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	hash := sha256.Sum256([]byte(token))
	return Session{
		TokenHash:         hash[:],
		UserID:            userID,
		CreatedAt:         now,
		LastSeenAt:        now,
		IdleExpiresAt:     now.Add(idleSessionLifetime),
		AbsoluteExpiresAt: now.Add(absoluteSessionLifetime),
	}, token, nil
}

func SessionTokenHash(token string) ([]byte, error) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(raw) != sessionTokenBytes {
		return nil, ErrUnauthenticated
	}
	hash := sha256.Sum256([]byte(token))
	return hash[:], nil
}

func RefreshSession(session Session, now time.Time) (Session, bool) {
	if !now.Before(session.IdleExpiresAt) || !now.Before(session.AbsoluteExpiresAt) {
		return session, false
	}
	session.LastSeenAt = now
	session.IdleExpiresAt = now.Add(idleSessionLifetime)
	if session.IdleExpiresAt.After(session.AbsoluteExpiresAt) {
		session.IdleExpiresAt = session.AbsoluteExpiresAt
	}
	return session, true
}
