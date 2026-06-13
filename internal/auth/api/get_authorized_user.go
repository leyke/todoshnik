package api

import (
	"context"

	appuser "todoshnik/internal/user"
)

func (h *Handler) GetAuthorizedUser(ctx context.Context, rawToken string) (*appuser.User, error) {
	token, err := h.tokenService.Get(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.GetById(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
