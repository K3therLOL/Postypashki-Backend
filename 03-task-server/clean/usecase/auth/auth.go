package usecase

import (
	"errors"
	domain "taskserver/clean/domain/auth"
)

var (
	ErrUserAlreadyExists   = errors.New("User already exists.")
	ErrUserNotExists       = errors.New("User doesn't exist. Please register first.")
	ErrWithPasswordHashing = errors.New("Couldn't hash password.")
	ErrWithUserSave        = errors.New("Couldn't save user.")
	ErrWithSessionCreating = errors.New("Couldn't create session.")
	ErrWithSessionSaving   = errors.New("Couldn't save session.")
	ErrWrongPassword       = errors.New("Wrong password.")
)

type AuthHandler struct {
	userRepo        domain.UserRepository
	sessionProvider domain.SessionProvider
	sessionRepo     domain.SessionRepository
	hasher          domain.Hasher
}

func NewAuthHandler(userRepo domain.UserRepository,
	sessionProvider domain.SessionProvider,
	sessionRepo domain.SessionRepository,
	hasher domain.Hasher) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, sessionProvider: sessionProvider, sessionRepo: sessionRepo, hasher: hasher}
}

func (handler *AuthHandler) Register(username, password string) error {
	_, ok := handler.userRepo.Exist(username)
	if !ok {
		return ErrUserAlreadyExists
	}

	passwordHash, err := handler.hasher.HashPassword(password)
	if err != nil {
		return ErrWithPasswordHashing
	}

	userObj := domain.NewUser(username, passwordHash)
	sessionObj, err := handler.sessionProvider.Create()
	if err != nil {
		return ErrWithSessionCreating
	}

	if err := handler.sessionRepo.Save(sessionObj); err != nil {
		return ErrWithSessionSaving
	}

	// UserID in own repo should be equal UserID in session
	userObj.ID = sessionObj.UserID
	if err := handler.userRepo.Save(userObj); err != nil {
		return ErrWithUserSave
	}

	return nil
}

func (handler *AuthHandler) Login(username, password string) error {
	userObj, ok := handler.userRepo.Exist(username)
	if !ok {
		return ErrUserNotExists
	}

	passwordHash, err := handler.hasher.HashPassword(password)
	if err != nil {
		return ErrWithPasswordHashing
	}
	if !handler.hasher.CheckPassword(passwordHash, password) {
		return ErrWrongPassword
	}

	sessionObj, err := handler.sessionProvider.Get(userObj.ID)
	if err == nil {
		handler.sessionProvider.Update(sessionObj)
	}
	return nil
}
