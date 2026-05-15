package bot

import "context"

func (h *Handler) DeleteTask(ctx context.Context, taskID string) error {
	response, err := h.api.Delete(ctx, "/tasks/"+taskID+"")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
