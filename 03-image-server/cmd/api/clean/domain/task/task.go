package domain

import (
	"github.com/google/uuid"
)

type TaskObject struct {
	TaskID uuid.UUID
	Status string
}

type TaskRepository interface {
	Save(uuid uuid.UUID) error
	Update(uuid uuid.UUID) error
	Get(uuid uuid.UUID) (*TaskObject, error)
}

type ResultRepository interface {
	Get(task_id string) (string, error)
	Set(task_id, imgUrl string) error
}
