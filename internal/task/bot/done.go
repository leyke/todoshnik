package bot

import (
	"context"
	"encoding/json"

	"todoshnik/internal/task"
)

func (h *Handler) DoneTask(ctx context.Context, taskID string) (string, error) {
	response, err := h.api.Post(ctx, taskAPIURL+taskID+"/done", nil)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var task task.Task

	err = json.NewDecoder(response.Body).Decode(&task)
	if err != nil {
		return "", err
	}

	return GetTaskRowText(task), nil
}
