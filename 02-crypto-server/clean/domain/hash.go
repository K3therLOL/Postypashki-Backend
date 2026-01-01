package domain

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hash, password string) bool
} 
