package repository

import (
	"cryptoserver/clean/domain"
)

type Rai struct {
	storage map[string]string // username -> hash
}

func NewRai() *Rai {
	return &Rai{storage: make(map[string]string)}
}

func (r *Rai) Save(user *domain.User) error {
	r.storage[user.Username] = user.PasswordHash
	return nil
}

func (r *Rai) Exist(username string) *domain.User {
	hash, ok := r.storage[username]
	if !ok {
		return nil
	}
	return &domain.User{Username: username, PasswordHash: hash}
}
