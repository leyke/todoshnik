package task

import (
	"context"
	"todoshnik/internal/domain"
)

type TaskRepositoryInterface interface {
	List(ctx context.Context, filter domain.TaskFilter) []*domain.Task
	GetByID(ctx context.Context, id int, scope domain.AccessScope) (*domain.Task, error)

	Create(ctx context.Context, task *domain.Task) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, task *domain.Task) error
}
