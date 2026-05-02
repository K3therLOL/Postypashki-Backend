package repository

import (
	"errors"
	domain "taskserver/clean/domain/task"

	"github.com/google/uuid"
)

const (
	InProgress = "in_progress"
	Ready      = "ready"
)

var (
	ErrNoId = errors.New("No object by your uuid.")
)

type TaskRepository struct {
	repo map[uuid.UUID]*domain.TaskObject // task_id -> task_object
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		repo: make(map[uuid.UUID]*domain.TaskObject),
	}
}

func (r *TaskRepository) Save(uuid uuid.UUID) error {
	r.repo[uuid] = &domain.TaskObject{
		TaskID: uuid,
		Status: InProgress,
	}
	return nil
}

func (r *TaskRepository) Update(uuid uuid.UUID) error {
	taskobj, ok := r.repo[uuid]
	if !ok {
		return ErrNoId
	}

	taskobj.Status = Ready
	return nil
}

func (r *TaskRepository) Get(uuid uuid.UUID) (*domain.TaskObject, error) {
	taskobj, ok := r.repo[uuid]
	if !ok {
		return nil, ErrNoId
	}

	return taskobj, nil
}
