package auth

import "context"

type Store interface {
	CreateUserWithSession(context.Context, User, Session) error
	FindUserByEmail(context.Context, string) (User, error)
	CreateSession(context.Context, Session) error
	FindSessionUser(context.Context, []byte) (SessionUser, error)
	TouchSession(context.Context, []byte, Session) (bool, error)
	DeleteSession(context.Context, []byte) error
	CreateAdmin(context.Context, User) error
}
