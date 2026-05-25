package bot

import (
	"context"
	"todoshnik/internal/task"
)

func (h *Handler) List(ctx context.Context, status string) ([]*task.Task, error) {
	response, err := h.taskClient.ListTasks(ctx, status)
	if err != nil {
		return nil, err
	}

	return response, nil
}
