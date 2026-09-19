package auth

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           string
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
	Role         Role
}

type Session struct {
	TokenHash         []byte
	UserID            string
	CreatedAt         time.Time
	LastSeenAt        time.Time
	IdleExpiresAt     time.Time
	AbsoluteExpiresAt time.Time
}

type SessionUser struct {
	Session Session
	User    User
}

type RegisterInput struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

type LoginInput struct {
	Email    string
	Password string
}

type AdminInput struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

type Authenticated struct {
	User         User
	SessionToken string
	ExpiresAt    time.Time
}
