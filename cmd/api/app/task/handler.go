package api

import (
	"context"
	"log"

	"todoshnik/internal/domains/task"
	"todoshnik/internal/domains/user"
	"todoshnik/internal/infrastructure/identity"
)

type UserGetter interface {
	GetById(ctx context.Context, userID int) (*user.User, error)
}

type TaskService interface {
	Add(ctx context.Context, title string, userID int) (*task.Task, error)
	List(ctx context.Context, filter task.TaskFilter) ([]*task.Task, error)
	Update(ctx context.Context, taskId int, title string, done bool, scope identity.AccessScope) (*task.Task, error)
	Get(ctx context.Context, taskId int, scope identity.AccessScope) (*task.Task, error)
	MarkDone(ctx context.Context, taskId int, scope identity.AccessScope) (*task.Task, error)
	Delete(ctx context.Context, taskId int, scope identity.AccessScope) error
}

type Handler struct {
	service    TaskService
	userGetter UserGetter
	logger     *log.Logger
}

func NewHandler(s TaskService, ug UserGetter, l *log.Logger) *Handler {
	return &Handler{service: s, userGetter: ug, logger: l}
}
