package bot

import (
	"context"
	"encoding/json"
	"todoshnik/internal/task"
)

func (h Handler) Add(ctx context.Context, taskTitle string) (*task.Task, error) {
	response, err := h.api.Post(ctx, "/tasks", CreateTaskRequest{
		Title: taskTitle,
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var task *task.Task

	err = json.NewDecoder(response.Body).Decode(&task)
	if err != nil {
		return nil, err
	}

	return task, nil
}
