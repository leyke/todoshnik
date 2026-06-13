package bot

import (
	"context"
	"encoding/json"
	"fmt"

	apierrors "todoshnik/internal/api/errors"
)

func (h *Handler) SignInUser(ctx context.Context, tgUser TgLoginRequestDto) (string, error) {
	response, err := h.api.Post(
		ctx,
		"/auth/tg/auto-reg",
		tgUser,
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseInfo UserAuthInfoResponseDto

	err = json.NewDecoder(response.Body).Decode(&responseInfo)
	if err != nil {
		return "", apierrors.ErrInvalidJSON
	}

	if responseInfo.AccessToken != "" {
		return responseInfo.AccessToken, nil
	}

	return "", fmt.Errorf(
		"Ошибка получения токена из: %v",
		responseInfo,
	)
}
