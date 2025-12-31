package composure

import (
	"cryptoserver/domain"
	"cryptoserver/usecase"
	"cryptoserver/repository"
	"cryptoserver/security"
	"cryptoserver/controller"
)

func NewAuth() *controller.Auth {
	repo := repository.NewRai()
	hasher := security.NewHasher()
}
