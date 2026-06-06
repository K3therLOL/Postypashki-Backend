package repository

import (
	"encoding/json"
	"errors"
	"fmt"
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
	r.repo[uuid] = taskobj
	return nil
}

func (r *TaskRepository) Get(uuid uuid.UUID) (*domain.TaskObject, error) {
	taskobj, ok := r.repo[uuid]
	if !ok {
		return nil, ErrNoId
	}

	return taskobj, nil
}

type brokerDTO struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

func (r *TaskRepository) Handle(body []byte) error {
	brokerObj := &brokerDTO{}
	if err := json.Unmarshal(body, brokerObj); err != nil {
		return err
	}

	taskID, err := uuid.Parse(brokerObj.TaskID)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := r.Update(taskID); err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("%s broker uuid updated\n", brokerObj.TaskID)
	return nil
}
