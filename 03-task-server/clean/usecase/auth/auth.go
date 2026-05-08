package usecase

import (
	"errors"
	"log"
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
	ErrAuthFirst           = errors.New("No session given. Please authenticate first.")
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

func (handler *AuthHandler) Register(username, password string) (string, error) {
	_, exist := handler.userRepo.Exist(username)
	if exist {
		return "", ErrUserAlreadyExists
	}

	passwordHash, err := handler.hasher.HashPassword(password)
	if err != nil {
		return "", ErrWithPasswordHashing
	}

	userObj := domain.NewUser(username, passwordHash)
	sessionObj, err := handler.sessionProvider.Create()
	if err != nil {
		return "", ErrWithSessionCreating
	}

	if err := handler.sessionRepo.Save(sessionObj); err != nil {
		return "", ErrWithSessionSaving
	}
	// UserID in own repo should be equal UserID in session
	userObj.ID = sessionObj.UserID
	if err := handler.userRepo.Save(userObj); err != nil {
		return "", ErrWithUserSave
	}

	return sessionObj.Token, nil
}

func (handler *AuthHandler) updateSession(session *domain.Session) error {
	if err := handler.sessionRepo.Delete(session); err != nil {
		return err
	}

	newSession, err := handler.sessionProvider.Create()
	if err != nil {
		return err
	}

	newSession.UserID = session.UserID
	if err := handler.sessionRepo.Save(newSession); err != nil {
		return err
	}

	return nil
}

func (handler *AuthHandler) Login(username, password string) (string, error) {
	userObj, ok := handler.userRepo.Exist(username)
	if !ok {
		return "", ErrUserNotExists
	}

	passwordHash, err := handler.hasher.HashPassword(password)
	if err != nil {
		return "", ErrWithPasswordHashing
	}
	if !handler.hasher.CheckPassword(passwordHash, password) {
		return "", ErrWrongPassword
	}

	sessionObj, ok := handler.sessionRepo.GetByUserID(userObj.ID)
	if !ok {
		handler.updateSession(sessionObj)
	}
	return sessionObj.Token, nil
}

func (handler *AuthHandler) GetSession(token string) (string, error) {
	sessionObj, ok := handler.sessionRepo.GetByToken(token)
	if !ok {
		return "", ErrAuthFirst
	}

	log.Printf("User %d has %s token that expires %v\n", sessionObj.UserID, sessionObj.Token, sessionObj.ExpiresAt)
	return sessionObj.Token, nil
}
