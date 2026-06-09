package cli

import (
	"context"
	"fmt"
)

func (h *Handler) clearTokens(ctx context.Context) {
	deleted, err := h.tokenService.ClearExpiredTokens(ctx)
	if err != nil {
		printErr(err, "Ошибка при удалении токенов")
		return
	}
	printInfo(fmt.Sprintf("Удалены токены: %v\n", deleted))
}
