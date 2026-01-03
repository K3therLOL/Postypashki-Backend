package security

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
}

func NewHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (h *BcryptHasher) CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
