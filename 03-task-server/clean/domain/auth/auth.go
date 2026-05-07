package domain

import "time"

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

func NewUser(username, passwordHash string) *User {
	return &User{Username: username, PasswordHash: passwordHash}
}

type UserRepository interface {
	Save(user *User) error
	Exist(username string) (*User, bool)
}

type Session struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time
}

type SessionProvider interface {
	Create() (*Session, error)
	Refresh() error
}

type SessionRepository interface {
	Save(session *Session) error
	Delete(session *Session) error
	GetByUserID(userID int) (*Session, bool)
	GetByToken(token string) (*Session, bool)
}

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hash, password string) bool
}
