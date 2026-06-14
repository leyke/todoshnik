package bot

import (
	"context"
	"encoding/json"
	"fmt"

	apierrors "todoshnik/internal/api/errors"
	tokenerror "todoshnik/internal/auth/token/errors"
)

const tgSignInURL string = "/auth/tg/auto-reg"

func (h *Handler) SignInUser(ctx context.Context, tgUser TgLoginRequestDto) (*AuthInfo, error) {
	response, err := h.api.Post(
		ctx,
		tgSignInURL,
		tgUser,
	)
	defer func() {
		_ = response.Body.Close()
	}()

	if err != nil {
		return nil, err
	}

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
