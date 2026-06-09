package bot

import (
	"context"
	"strconv"

	authbot "todoshnik/internal/auth/bot"
	authcontext "todoshnik/internal/auth/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleAuth(ctx context.Context, user *tgbotapi.User) context.Context {
	tokenCacheKey := "user:" + strconv.FormatInt(user.ID, 10) + ":tg-api-token-key"

	// попытка забрать из кеша
	token, _ := h.cache.Get(ctx, tokenCacheKey).Result()

	if token != "" {
		ctx = authcontext.SetToken(ctx, token)
		return ctx
	}

	// попытка сгенерировать через апи
	token, err := h.authHandler.GetToken(ctx, authbot.TgLoginRequestDto{
		TgUserID: user.ID,
		Name:     user.UserName,
	})

	if err != nil {
		return ctx
	}

	// вставить в кеш
	h.cache.Set(ctx, tokenCacheKey, token, botTokenTtl)

	// вставить в контекст
	ctx = authcontext.SetToken(ctx, token)

	return ctx
}

func (h *Handler) handleWelcome(ctx context.Context, user *tgbotapi.User) error {
	tokenCacheKey := "user:" + strconv.FormatInt(user.ID, 10) + ":tg-api-token-key"

	token, err := h.authHandler.SignInUser(ctx, authbot.TgLoginRequestDto{
		TgUserID: user.ID,
		Name:     user.UserName,
	})
	if err != nil {
		return err
	}

	h.cache.Set(ctx, tokenCacheKey, token, botTokenTtl)

	return nil
}
