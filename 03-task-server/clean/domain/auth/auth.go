package domain

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
	Exist(username string) *User
}

type Session struct {
	userID    int
	sessionID string
}

func NewSession() *Session {
	return &Session{}
}

type SessionProvider interface {
	Create() (*Session, error)
	Get(sid string) (*Session, error)
	Delete(sid string) error
	Refresh() error
}

type SessionRepository interface {
	Save(session *Session) error
	Exist(sid string) *Session
}

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hash, password string) bool
}
