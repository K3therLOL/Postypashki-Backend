package composure

import (
	"os"
	controller "taskserver/clean/controller/auth"
	usecase "taskserver/clean/usecase/auth"
	hasher "taskserver/hash"
	provider "taskserver/provider/session"
	session "taskserver/repository/auth/sessions/database"
	user "taskserver/repository/auth/users/database"
)

func NewAuthController() *controller.AuthController {
	connString := os.Getenv("DB_CONN_STRING")

	userRepo := user.NewUserRepository(connString)
	sessionRepo := session.NewSessionRepository(connString)
	sessionProvider := provider.NewSessionProvider(24)
	hasher := hasher.NewHasher(10)
	usecase := usecase.NewAuthHandler(userRepo, sessionProvider, sessionRepo, hasher)
	return controller.NewAuthController(usecase)
}
