package composure

import (
	"cryptoserver/clean/usecase"
	"cryptoserver/repository"
	"cryptoserver/security"
	"cryptoserver/clean/controller"
)

func NewAuth() *controller.Auth {
	repo := repository.NewRai()
	hasher := security.NewHasher()
	usecase := usecase.NewAuth(repo, hasher)
	auth := controller.NewAuth(usecase)
	return auth
}
