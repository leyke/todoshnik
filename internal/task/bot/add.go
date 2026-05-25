package bot

import (
	"context"
	"todoshnik/internal/task"
)

func (h *Handler) Add(ctx context.Context, taskTitle string) (*task.Task, error) {
	response, err := h.taskClient.CreateTask(ctx, taskTitle)
	if err != nil {
		return nil, err
	}

	return response, nil
}
