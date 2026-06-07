package api

import (
	"context"
	"todoshnik/internal/user"
)

func (h *Handler) GetAuthorizedUser(ctx context.Context, rawToken string) (*user.User, error) {
	token, err := h.tokenService.Get(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	// линтер любезно подсказывает что имя переменной конфликтует и перетирает название импортируемого пакета
	user, err := h.userService.Get(ctx, token.UserID, "")
	if err != nil {
		return nil, err
	}

	return user, nil
}
