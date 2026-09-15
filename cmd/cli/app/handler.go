package cli

import (
	"context"
	"fmt"

	"todoshnik/internal/domains/task"
	"todoshnik/internal/infrastructure/identity"
)

type TaskService interface {
	Add(ctx context.Context, title string, taskID int) (*task.Task, error)
	List(ctx context.Context, filter task.TaskFilter) ([]*task.Task, error)
	MarkDone(ctx context.Context, taskID int, scope identity.AccessScope) (*task.Task, error)
	Delete(ctx context.Context, taskID int, scope identity.AccessScope) error
}

type TokenService interface {
	ClearExpiredTokens(ctx context.Context) (int, error)
}

type Handler struct {
	taskService  TaskService
	tokenService TokenService
}

func NewHandler(srvTask TaskService, srvToken TokenService) *Handler {
	return &Handler{
		taskService:  srvTask,
		tokenService: srvToken,
	}
}

func (h *Handler) Run(args []string) {
	if len(args) < 2 {
		printHelp()
		return
	}

	ctx := context.Background()

	switch args[1] {
	case "add":
		h.add(ctx, args)
	case "list":
		h.list(ctx, args)
	case "done":
		h.done(ctx, args)
	case "delete":
		h.delete(ctx, args)
	case "clear-tokens":
		h.clearTokens(ctx)
	default:
		printErr(fmt.Errorf("неизвестная команда: %s", args[1]), "")
	}
}
