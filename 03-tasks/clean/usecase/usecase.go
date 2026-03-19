package usecase

import "task/clean/domain"

type Repository struct {
	repo domain.TaskRepository
}

func NewRepository(repo domain.TaskRepository) *Repository {
	return &Repository{repo: repo}
}
