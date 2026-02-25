package composure

import (
	"cryptoserver/clean/controller"
	"cryptoserver/clean/usecase"
	repository "cryptoserver/repository/database"
	"cryptoserver/security"
)

func NewAuth() *controller.Auth {
	connString := "postgres://kether:abc@localhost:5432/postgres?sslmode=disable"
	repo := repository.NewStorage(connString)
	hasher := security.NewHasher()
	usecase := usecase.NewAuth(repo, hasher)
	auth := controller.NewAuth(usecase)
	return auth
}
