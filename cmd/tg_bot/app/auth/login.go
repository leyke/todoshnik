package auth

import (
	"context"
	"encoding/json"
	"fmt"

	apierrors "todoshnik/internal/api/errors"
	tokenerror "todoshnik/internal/domains/token/errors"
)

const tgAuthLoginURL = "/auth/tg/login"

func (h *Handler) Login(ctx context.Context, tgUser TgLoginRequestDto) (*AuthInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, h.requestTimeout)
	defer cancel()

	response, err := h.api.Post(
		ctx,
		tgAuthLoginURL,
		tgUser,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var responseInfo UserAuthInfoResponseDto
	if err := json.NewDecoder(response.Body).Decode(&responseInfo); err != nil {
		return nil, fmt.Errorf("%w: %v", apierrors.ErrInvalidJSON, err)
	}

	if responseInfo.AccessToken == "" {
		return nil, fmt.Errorf(
			"%w for user id %d",
			tokenerror.ErrNotFound,
			responseInfo.UserID,
		)
	}

	return &AuthInfo{
		UserID:      responseInfo.UserID,
		AccessToken: responseInfo.AccessToken,
	}, nil
}
