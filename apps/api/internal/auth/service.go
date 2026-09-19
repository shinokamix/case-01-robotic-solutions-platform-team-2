package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	store     Store
	hasher    PasswordHasher
	now       func() time.Time
	dummyHash string
}

func NewService(store Store, now func() time.Time) (*Service, error) {
	if now == nil {
		now = time.Now
	}
	hasher := PasswordHasher{}
	dummyHash, err := hasher.Hash("invalid-credential-padding")
	if err != nil {
		return nil, fmt.Errorf("prepare credential comparison: %w", err)
	}
	return &Service{store: store, hasher: hasher, now: now, dummyHash: dummyHash}, nil
}

func (service *Service) Register(ctx context.Context, input RegisterInput) (Authenticated, error) {
	email, err := CanonicalEmail(input.Email)
	if err != nil {
		return Authenticated{}, err
	}
	firstName, err := NormalizeFirstName(input.FirstName)
	if err != nil {
		return Authenticated{}, err
	}
	lastName, err := NormalizeLastName(input.LastName)
	if err != nil {
		return Authenticated{}, err
	}
	passwordHash, err := service.hasher.Hash(input.Password)
	if err != nil {
		return Authenticated{}, err
	}

	user := User{
		ID:           ulid.Make().String(),
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		Role:         RoleUser,
	}
	now := service.now().UTC()
	session, token, err := NewSession(user.ID, now)
	if err != nil {
		return Authenticated{}, err
	}
	if err := service.store.CreateUserWithSession(ctx, user, session); err != nil {
		return Authenticated{}, err
	}
	return Authenticated{User: user, SessionToken: token, ExpiresAt: session.AbsoluteExpiresAt}, nil
}

func (service *Service) Login(ctx context.Context, input LoginInput) (Authenticated, error) {
	email, err := CanonicalEmail(input.Email)
	if err != nil {
		service.hasher.Compare(service.dummyHash, input.Password)
		return Authenticated{}, ErrInvalidCredentials
	}
	user, err := service.store.FindUserByEmail(ctx, email)
	if err != nil {
		service.hasher.Compare(service.dummyHash, input.Password)
		if errors.Is(err, ErrNotFound) {
			return Authenticated{}, ErrInvalidCredentials
		}
		return Authenticated{}, err
	}
	if !service.hasher.Compare(user.PasswordHash, input.Password) {
		return Authenticated{}, ErrInvalidCredentials
	}

	now := service.now().UTC()
	session, token, err := NewSession(user.ID, now)
	if err != nil {
		return Authenticated{}, err
	}
	if err := service.store.CreateSession(ctx, session); err != nil {
		return Authenticated{}, err
	}
	return Authenticated{User: user, SessionToken: token, ExpiresAt: session.AbsoluteExpiresAt}, nil
}

func (service *Service) CurrentUser(ctx context.Context, token string) (User, error) {
	tokenHash, err := SessionTokenHash(token)
	if err != nil {
		return User{}, ErrUnauthenticated
	}
	result, err := service.store.FindSessionUser(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, ErrUnauthenticated
		}
		return User{}, err
	}

	now := service.now().UTC()
	refreshed, valid := RefreshSession(result.Session, now)
	if !valid {
		_ = service.store.DeleteSession(ctx, tokenHash)
		return User{}, ErrUnauthenticated
	}
	updated, err := service.store.TouchSession(ctx, tokenHash, refreshed)
	if err != nil {
		return User{}, err
	}
	if !updated {
		return User{}, ErrUnauthenticated
	}
	return result.User, nil
}

func (service *Service) CreateAdmin(ctx context.Context, input AdminInput) (User, error) {
	email, err := CanonicalEmail(input.Email)
	if err != nil {
		return User{}, err
	}
	firstName, err := NormalizeFirstName(input.FirstName)
	if err != nil {
		return User{}, err
	}
	lastName, err := NormalizeLastName(input.LastName)
	if err != nil {
		return User{}, err
	}
	passwordHash, err := service.hasher.Hash(input.Password)
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:           ulid.Make().String(),
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: passwordHash,
		Role:         RoleAdmin,
	}
	if err := service.store.CreateAdmin(ctx, user); err != nil {
		return User{}, err
	}
	return user, nil
}
