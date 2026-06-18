package api

import (
	"context"

	apierror "todoshnik/internal/api/errors"
)

func (h *Handler) Authorize(
	ctx context.Context,
	rawToken string,
) (int, error) {
	if rawToken == "" {
		return 0, apierror.ErrUnauth
	}

	token, err := h.tokenService.Get(ctx, rawToken)
	if err != nil {
		return 0, err
	}

	_, err = h.userService.GetById(ctx, token.UserID)
	if err != nil {
		return 0, err
	}

	return token.UserID, nil
}
