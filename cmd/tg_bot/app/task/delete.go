package bot

import "context"

func (h *Handler) DeleteTask(ctx context.Context, taskID string) error {
	response, err := h.api.Delete(ctx, taskAPIURL+taskID+"")
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	return nil
}
