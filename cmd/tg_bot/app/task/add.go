package bot

import (
	"context"
	"encoding/json"

	"todoshnik/internal/domains/task"
)

func (h *Handler) Add(ctx context.Context, taskTitle string) (*task.Task, error) {
	response, err := h.api.Post(ctx, taskAPIURL, CreateTaskRequest{
		Title: taskTitle,
	})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	var t task.Task

	if err := json.NewDecoder(response.Body).Decode(&t); err != nil {
		return nil, err
	}

	return &t, nil
}
