package auth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type memoryStore struct {
	mu       sync.Mutex
	users    map[string]User
	sessions map[string]Session
}

func newMemoryStore() *memoryStore {
	return &memoryStore{users: make(map[string]User), sessions: make(map[string]Session)}
}

func (store *memoryStore) CreateUserWithSession(_ context.Context, user User, session Session) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	if _, exists := store.users[user.Email]; exists {
		return ErrEmailAlreadyExists
	}
	store.users[user.Email] = user
	store.sessions[string(session.TokenHash)] = session
	return nil
}

func (store *memoryStore) FindUserByEmail(_ context.Context, email string) (User, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	user, exists := store.users[email]
	if !exists {
		return User{}, ErrNotFound
	}
	return user, nil
}

func (store *memoryStore) CreateSession(_ context.Context, session Session) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.sessions[string(session.TokenHash)] = session
	return nil
}

func (store *memoryStore) FindSessionUser(_ context.Context, tokenHash []byte) (SessionUser, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	session, exists := store.sessions[string(tokenHash)]
	if !exists {
		return SessionUser{}, ErrNotFound
	}
	for _, user := range store.users {
		if user.ID == session.UserID {
			return SessionUser{Session: session, User: user}, nil
		}
	}
	return SessionUser{}, ErrNotFound
}

func (store *memoryStore) TouchSession(_ context.Context, tokenHash []byte, session Session) (bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if _, exists := store.sessions[string(tokenHash)]; !exists {
		return false, nil
	}
	store.sessions[string(tokenHash)] = session
	return true, nil
}

func (store *memoryStore) DeleteSession(_ context.Context, tokenHash []byte) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.sessions, string(tokenHash))
	return nil
}

func (store *memoryStore) CreateAdmin(_ context.Context, user User) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	if _, exists := store.users[user.Email]; exists {
		return ErrEmailAlreadyExists
	}
	store.users[user.Email] = user
	return nil
}

func TestServiceRegistrationLoginAndCurrentUser(t *testing.T) {
	store := newMemoryStore()
	now := time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)
	service, err := NewService(store, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}

	registered, err := service.Register(context.Background(), RegisterInput{
		Email: " USER+demo@EXAMPLE.COM ", FirstName: " Иван ", LastName: " Петров ", Password: "very-secure-password",
	})
	if err != nil {
		t.Fatal(err)
	}
	if registered.User.Email != "user+demo@example.com" || registered.User.Role != RoleUser {
		t.Fatalf("unexpected registered user: %+v", registered.User)
	}
	if registered.User.PasswordHash == "very-secure-password" || registered.User.PasswordHash == "" {
		t.Fatal("password must be stored only as a hash")
	}

	if _, err := service.Register(context.Background(), RegisterInput{
		Email: "user+demo@example.com", FirstName: "Иван", LastName: "Петров", Password: "very-secure-password",
	}); !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("duplicate registration error = %v", err)
	}

	if _, err := service.Login(context.Background(), LoginInput{Email: registered.User.Email, Password: "wrong-password"}); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("invalid login error = %v", err)
	}
	loggedIn, err := service.Login(context.Background(), LoginInput{Email: "USER+DEMO@EXAMPLE.COM", Password: "very-secure-password"})
	if err != nil {
		t.Fatal(err)
	}
	current, err := service.CurrentUser(context.Background(), loggedIn.SessionToken)
	if err != nil {
		t.Fatal(err)
	}
	if current.FirstName != "Иван" || current.LastName != "Петров" || current.Role != RoleUser {
		t.Fatalf("unexpected current user: %+v", current)
	}
}

func TestServiceRejectsExpiredSession(t *testing.T) {
	store := newMemoryStore()
	now := time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)
	service, err := NewService(store, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	registered, err := service.Register(context.Background(), RegisterInput{
		Email: "user@example.com", FirstName: "Иван", LastName: "Петров", Password: "very-secure-password",
	})
	if err != nil {
		t.Fatal(err)
	}
	now = now.Add(14 * 24 * time.Hour)
	if _, err := service.CurrentUser(context.Background(), registered.SessionToken); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expired session error = %v", err)
	}
}
