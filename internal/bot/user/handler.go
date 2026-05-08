package user

import (
	"context"
	"encoding/json"
	"fmt"
	"todoshnik/internal/client"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	api *client.ApiClient
}

func NewHandler(api *client.ApiClient) *Handler {
	return &Handler{api: api}
}

func (h Handler) GetToken(ctx context.Context, tgUser *tgbotapi.User) (string, error) {
	response, err := h.api.Post(
		ctx,
		"/auth/tg/login",
		TgLoginRequestDto{
			TgUserID: tgUser.ID,
			Name:     tgUser.UserName,
		},
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseInfo UserAuthInfoResponseDto

	err = json.NewDecoder(response.Body).Decode(&responseInfo)
	if err != nil {
		return "", err
	}

	if responseInfo.AccessToken != "" {
		return responseInfo.AccessToken, nil
	}

	return "", fmt.Errorf(
		"Ошибка получения токена из: %v",
		responseInfo,
	)
}

func (h Handler) SignInUser(ctx context.Context, tgUser *tgbotapi.User) (string, error) {
	response, err := h.api.Post(
		ctx,
		"/auth/tg/auto-reg",
		TgLoginRequestDto{
			TgUserID: tgUser.ID,
			Name:     tgUser.UserName,
		},
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseInfo UserAuthInfoResponseDto

	err = json.NewDecoder(response.Body).Decode(&responseInfo)
	if err != nil {
		return "", err
	}

	if responseInfo.AccessToken != "" {
		return responseInfo.AccessToken, nil
	}

	return "", fmt.Errorf(
		"Ошибка получения токена из: %v",
		responseInfo,
	)
}
