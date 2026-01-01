package usecase

import (
	"cryptoserver/clean/domain"
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("User already exists")
	ErrUserNotExists = errors.New("Username doesn't exist. Please register first.")
	ErrWrongPassword = errors.New("Wrong password.")
)

type Auth struct {
	ur domain.UserRepository
	h  domain.Hasher 
}

func NewAuth(ur domain.UserRepository, h domain.Hasher) *Auth {
	return &Auth{ur: ur, h: h}
}

func (usecase *Auth) Register(username, password string) error {
	if user := usecase.ur.Exist(username); user != nil {
		return ErrUserAlreadyExists
	}

	hash, err := usecase.h.HashPassword(password)
	if err != nil {
		return err
	}

	user := domain.NewUser(username, hash)
	if err := usecase.ur.Save(user); err != nil {
		return err
	}
	
	return nil
}

func (usecase *Auth) Login(username, password string) (*domain.User, error) {
	user := usecase.ur.Exist(username)
	if user == nil {
		return nil, ErrUserNotExists
	}

	hash := user.PasswordHash
	if !usecase.h.CheckPassword(hash, password) {
		return nil, ErrWrongPassword
	}

	return user, nil
}
