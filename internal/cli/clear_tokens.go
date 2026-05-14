package cli

import (
	"context"
	"fmt"
)

func (h *Handler) clearTokens(ctx context.Context) {
	deleted := h.tokenService.ClearExpiredTokens(ctx)
	printInfo(fmt.Sprintf("Удалены токены: %v\n", deleted))
}
