package task

import (
	"context"
	"todoshnik/internal/identity"
)

type Repository interface {
	List(ctx context.Context, filter TaskFilter) ([]*Task, error)
	GetByID(ctx context.Context, id int, scope identity.AccessScope) (*Task, error)

	Create(ctx context.Context, task *Task) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, task *Task) error
}
