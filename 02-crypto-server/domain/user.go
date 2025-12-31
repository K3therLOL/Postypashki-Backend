package domain

type User struct {
	Username string
	PasswordHash string
}

func NewUser(username, passwordHash string) *User {
	return &User{ Username: username, PasswordHash: passwordHash }
}

type UserRepository interface {
	Save(user *User) error
	Exist(username string) *User
}
