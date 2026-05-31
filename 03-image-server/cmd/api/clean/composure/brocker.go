package composure

import (
	"os"
	brocker "taskserver/brocker/rabbitmq/0.9"
	usecase "taskserver/clean/usecase/brocker"
)

func NewBrockerInteractor() *usecase.BrockerInteractor {
	brokerURI := os.Getenv("BROKER_URI")
	provider := brocker.NewRabbitProducer(brokerURI)
	interactor := usecase.NewBrockerInteractor(provider)
	return interactor
}
