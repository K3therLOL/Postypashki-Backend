package domain

type User struct {
	ID           int
	Login        string
	PasswordHash string
}

type UserRepository interface {
	Save(user *User) error
	Exist(login string) *User
}

type Session struct {
	userID    int
	sessionID string
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
