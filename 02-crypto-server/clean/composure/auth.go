package composure

import (
	"cryptoserver/usecase"
	"cryptoserver/repository"
	"cryptoserver/security"
	"cryptoserver/controller"
)

func NewAuth() *controller.Auth {
	repo := repository.NewRai()
	hasher := security.NewHasher()
	usecase := usecase.NewAuth(repo, hasher)
	auth := controller.NewAuth(usecase)
	return auth
}
