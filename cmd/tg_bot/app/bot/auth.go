package bot

import (
	"context"
	"fmt"

	"todoshnik/cmd/tg_bot/app/auth"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleAuth(ctx context.Context, user *tgbotapi.User) (string, error) {
	// попытка забрать из кеша
	token, err := h.tokenStorage.Get(ctx, user.ID)
	if err == nil && token != "" {
		return token, nil
	}

	// попытка сгенерировать через апи
	tokenInfo, err := h.authHandler.Login(ctx, auth.TgLoginRequestDto{
		TgUserID: user.ID,
		Name:     user.UserName,
	})

	if err != nil {
		return "", fmt.Errorf("ошибка авторизации: %w", err)
	}

	// вставить в кеш
	err = h.tokenStorage.Set(ctx, user.ID, token)
	if err != nil {
		return "", fmt.Errorf("ошибка запоминания токена: %w", err)
	}

	return tokenInfo.AccessToken, nil
}

func (h *Handler) handleWelcome(ctx context.Context, user *tgbotapi.User) error {
	tokenInfo, err := h.authHandler.SignInUser(ctx, auth.TgLoginRequestDto{
		TgUserID: user.ID,
		Name:     user.UserName,
	})
	if err != nil {
		return err
	}
	// вставить в кеш
	err = h.tokenStorage.Set(ctx, user.ID, tokenInfo.AccessToken)
	if err != nil {
		return err
	}

	return nil
}
