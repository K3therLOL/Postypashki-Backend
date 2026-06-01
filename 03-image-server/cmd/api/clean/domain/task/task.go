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
	Get(uuid string) (string, error)
	Set(uuid, imgUrl string) error
}
