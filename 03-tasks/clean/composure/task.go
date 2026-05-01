package composure

import (
	"taskserver/clean/usecase"
	resRepository "taskserver/repository/result/redis"
	taskRepository "taskserver/repository/task/rai"
)

func NewTaskInteractor() *usecase.TaskInteractor {
	taskRepo := taskRepository.NewTaskRepository()
	resRepo := resRepository.NewResultRepository(24)
	interactor := usecase.NewTaskInteractor(taskRepo, resRepo)
	return interactor
}
