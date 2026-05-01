package usecase

import (
	domain "taskserver/clean/domain/task"

	"github.com/google/uuid"
)

type TaskInteractor struct {
	taskRepo domain.TaskRepository
	resRepo  domain.ResultRepository
}

func NewTaskInteractor(taskRepo domain.TaskRepository, resRepo domain.ResultRepository) *TaskInteractor {
	return &TaskInteractor{
		taskRepo: taskRepo,
		resRepo:  resRepo,
	}
}

func (interactor *TaskInteractor) SaveTask(taskID uuid.UUID) error {
	return interactor.taskRepo.Save(taskID)
}

func (interactor *TaskInteractor) GetTaskStatus(taskID uuid.UUID) (string, error) {
	taskobj, err := interactor.taskRepo.Get(taskID)
	if err != nil {
		return "", err
	}

	return taskobj.Status, nil
}

func (interactor *TaskInteractor) UpdateTaskStatus(taskID uuid.UUID) error {
	return interactor.taskRepo.Update(taskID)
}

func (interactor *TaskInteractor) SaveResult(taskID, result string) error {
	return interactor.resRepo.Set(taskID, result)
}

func (interactor *TaskInteractor) GetResult(taskID string) (string, error) {
	return interactor.resRepo.Get(taskID)
}
