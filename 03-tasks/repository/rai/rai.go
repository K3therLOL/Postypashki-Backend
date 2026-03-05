package repository

import (
	"errors"
	"task/clean/domain"

	"github.com/google/uuid"
)

const (
	InProgress = "in_progress"
	Ready      = "ready"
)

var (
	ErrNoId = errors.New("No object by your uuid.")
)

type Repository struct {
	repo map[uuid.UUID]*domain.TaskObject // task_id -> task_object
}

func NewRepository() *Repository {
	return &Repository{
		repo: make(map[uuid.UUID]*domain.TaskObject),
	}
}

func (r *Repository) Save(uuid uuid.UUID) error {
	r.repo[uuid] = &domain.TaskObject{
		TaskID: uuid,
		Status: InProgress,
	}
	return nil
}

func (r *Repository) Update(uuid uuid.UUID) error {
	taskobj, ok := r.repo[uuid]
	if !ok {
		return ErrNoId
	}

	taskobj.Status = Ready
	return nil
}

func (r *Repository) Get(uuid uuid.UUID) (*domain.TaskObject, error) {
	taskobj, ok := r.repo[uuid]
	if !ok {
		return nil, ErrNoId
	}

	return taskobj, nil
}
