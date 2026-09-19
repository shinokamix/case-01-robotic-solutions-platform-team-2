package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth"
	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/postgres/sqlc"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (store *Store) CreateUserWithSession(ctx context.Context, user auth.User, session auth.Session) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	queries := sqlc.New(tx)
	if _, err := queries.CreateUser(ctx, createUserParams(user)); err != nil {
		return mapUserWriteError(err, auth.ErrEmailAlreadyExists)
	}
	if err := queries.CreateSession(ctx, createSessionParams(session)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (store *Store) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {
	user, err := sqlc.New(store.pool).GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return auth.User{}, auth.ErrNotFound
	}
	if err != nil {
		return auth.User{}, err
	}
	return mapUser(user), nil
}

func (store *Store) CreateSession(ctx context.Context, session auth.Session) error {
	return sqlc.New(store.pool).CreateSession(ctx, createSessionParams(session))
}

func (store *Store) FindSessionUser(ctx context.Context, tokenHash []byte) (auth.SessionUser, error) {
	row, err := sqlc.New(store.pool).GetSessionUser(ctx, tokenHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return auth.SessionUser{}, auth.ErrNotFound
	}
	if err != nil {
		return auth.SessionUser{}, err
	}
	return auth.SessionUser{
		Session: auth.Session{
			TokenHash:         row.TokenHash,
			UserID:            row.ID,
			CreatedAt:         row.SessionCreatedAt.Time,
			LastSeenAt:        row.LastSeenAt.Time,
			IdleExpiresAt:     row.IdleExpiresAt.Time,
			AbsoluteExpiresAt: row.AbsoluteExpiresAt.Time,
		},
		User: auth.User{
			ID:           row.ID,
			Email:        row.Email,
			FirstName:    row.FirstName,
			LastName:     row.LastName,
			PasswordHash: row.PasswordHash,
			Role:         auth.Role(row.Role),
		},
	}, nil
}

func (store *Store) TouchSession(ctx context.Context, tokenHash []byte, session auth.Session) (bool, error) {
	count, err := sqlc.New(store.pool).TouchSession(ctx, sqlc.TouchSessionParams{
		TokenHash:     tokenHash,
		LastSeenAt:    timestamp(session.LastSeenAt),
		IdleExpiresAt: timestamp(session.IdleExpiresAt),
	})
	return count == 1, err
}

func (store *Store) DeleteSession(ctx context.Context, tokenHash []byte) error {
	return sqlc.New(store.pool).DeleteSession(ctx, tokenHash)
}

func (store *Store) CreateAdmin(ctx context.Context, user auth.User) error {
	queries := sqlc.New(store.pool)
	_, err := queries.CreateUser(ctx, createUserParams(user))
	if isAdminConflict(err) {
		return auth.ErrAdminAlreadyExists
	}
	if !isEmailConflict(err) {
		return err
	}
	existing, findErr := queries.GetUserByEmail(ctx, user.Email)
	if findErr != nil {
		return findErr
	}
	if auth.Role(existing.Role) == auth.RoleAdmin {
		return auth.ErrAdminAlreadyExists
	}
	return auth.ErrEmailAlreadyExists
}

func createUserParams(user auth.User) sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PasswordHash: user.PasswordHash,
		Role:         string(user.Role),
	}
}

func createSessionParams(session auth.Session) sqlc.CreateSessionParams {
	return sqlc.CreateSessionParams{
		TokenHash:         session.TokenHash,
		UserID:            session.UserID,
		CreatedAt:         timestamp(session.CreatedAt),
		LastSeenAt:        timestamp(session.LastSeenAt),
		IdleExpiresAt:     timestamp(session.IdleExpiresAt),
		AbsoluteExpiresAt: timestamp(session.AbsoluteExpiresAt),
	}
}

func mapUser(user sqlc.User) auth.User {
	return auth.User{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PasswordHash: user.PasswordHash,
		Role:         auth.Role(user.Role),
	}
}

func timestamp(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: value, Valid: true}
}

func mapUserWriteError(err error, conflict error) error {
	if isEmailConflict(err) {
		return conflict
	}
	return err
}

func isEmailConflict(err error) bool {
	var postgresError *pgconn.PgError
	return errors.As(err, &postgresError) && postgresError.Code == "23505" && postgresError.ConstraintName == "users_email_key"
}

func isAdminConflict(err error) bool {
	var postgresError *pgconn.PgError
	return errors.As(err, &postgresError) && postgresError.Code == "23505" && postgresError.ConstraintName == "users_single_admin_idx"
}

var _ auth.Store = (*Store)(nil)
