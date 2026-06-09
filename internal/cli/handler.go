package cli

import (
	"context"
	"fmt"

	"todoshnik/internal/app"
	"todoshnik/internal/auth/token"
	"todoshnik/internal/task"
)

type Handler struct {
	taskService  *task.Service
	tokenService *token.Service
}

func NewHandler(services *app.Services) *Handler {
	return &Handler{
		taskService:  services.TaskService,
		tokenService: services.TokenService,
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
		printErr(fmt.Errorf("Неизвестная команда: %s\n", args[1]), "")
	}
}
