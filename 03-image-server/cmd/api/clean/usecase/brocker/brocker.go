package usecase

import domain "taskserver/clean/domain/brocker"

type BrockerInteractor struct {
	provider domain.BrockerProvider
}

func NewBrockerInteractor(provider domain.BrockerProvider) *BrockerInteractor {
	return &BrockerInteractor{provider: provider}
}

func (interactor *BrockerInteractor) Send(message string) error {
	return interactor.provider.Send(message)
}
