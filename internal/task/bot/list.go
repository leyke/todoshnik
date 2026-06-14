package bot

import (
	"context"
	"encoding/json"
	"net/url"

	"todoshnik/internal/task"
)

func (h *Handler) List(ctx context.Context, status string) ([]*task.Task, error) {
	params := url.Values{}
	if status != "" {
		params.Add("status", status)
	}

	response, err := h.api.Get(ctx, taskAPIURL, params)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var tasks []*task.Task

	err = json.NewDecoder(response.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
