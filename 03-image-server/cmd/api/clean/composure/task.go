package composure

import (
	"os"
	usecase "taskserver/clean/usecase/task"
	resRepository "taskserver/repository/result/postgres"
	taskRepository "taskserver/repository/task/rai"
)

func NewTaskInteractor() *usecase.TaskInteractor {
	taskRepo := taskRepository.NewTaskRepository()
	connString := os.Getenv("DB_CONN_STRING")
	resRepo := resRepository.NewImageRepository(connString)
	interactor := usecase.NewTaskInteractor(taskRepo, resRepo)
	return interactor
}
