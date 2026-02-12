package composure

import (
	"cryptoserver/clean/controller"
	"cryptoserver/clean/usecase"
	"cryptoserver/repository"
	"cryptoserver/security"
)

func NewAuth() *controller.Auth {
	repo := repository.NewRai()
	hasher := security.NewHasher()
	usecase := usecase.NewAuth(repo, hasher)
	auth := controller.NewAuth(usecase)
	return auth
}
