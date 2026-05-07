package composure

import (
	"os"
	usecase "taskserver/clean/usecase/auth"
	hasher "taskserver/hash"
	provider "taskserver/provider/session"
	session "taskserver/repository/auth/sessions/database"
	user "taskserver/repository/auth/users/database"
)

func NewAuthHandler() *usecase.AuthHandler {
	connString := os.Getenv("DB_CONN_STRING")

	userRepo := user.NewUserRepository(connString)
	sessionRepo := session.NewSessionRepository(connString)
	sessionProvider := provider.NewSessionProvider(24)
	hasher := hasher.NewHasher(10)
	return usecase.NewAuthHandler(userRepo, sessionProvider, sessionRepo, hasher)
}
