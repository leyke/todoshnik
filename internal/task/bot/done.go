package bot

import (
	"context"
	"encoding/json"
	"todoshnik/internal/task"
)

func (h Handler) DoneTask(ctx context.Context, taskID string) (string, error) {
	response, err := h.api.Post(ctx, "/tasks/"+taskID+"/done", nil)
	if err != nil {
		return "", err
	}

	var task task.Task

	err = json.NewDecoder(response.Body).Decode(&task)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return getTaskRowText(task), nil
}
