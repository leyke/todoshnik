package task

import "todoshnik/internal/domain"

type TaskRepositoryInreface interface {
	List(filter domain.TaskFilter) []*domain.Task
	GetByID(id int, scope domain.AccessScope) (*domain.Task, error)

	Create(task *domain.Task) (*domain.Task, error)
	Update(task *domain.Task) error
	Delete(task *domain.Task) error
}
