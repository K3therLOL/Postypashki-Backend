package composure

import (
	"os"
	brocker "taskserver/brocker/rabbitmq/0.9"
	brokerUsecase "taskserver/clean/usecase/brocker"
	taskUsecase "taskserver/clean/usecase/task"
	resRepository "taskserver/repository/result/postgres"
	taskRepository "taskserver/repository/task/rai"
)

var taskRepo *taskRepository.TaskRepository

func init() {
	taskRepo = taskRepository.NewTaskRepository()
}

func NewBrockerInteractor() *brokerUsecase.BrockerInteractor {
	brokerURI := os.Getenv("BROKER_URI")
	provider := brocker.NewRabbitProducer(brokerURI, taskRepo)
	interactor := brokerUsecase.NewBrockerInteractor(provider)
	return interactor
}

func NewTaskInteractor() *taskUsecase.TaskInteractor {
	connString := os.Getenv("DB_CONN_STRING")
	resRepo := resRepository.NewImageRepository(connString)
	interactor := taskUsecase.NewTaskInteractor(taskRepo, resRepo)
	return interactor
}
