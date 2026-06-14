package bot

import (
	"context"
	"fmt"

	authbot "todoshnik/internal/auth/bot"
	authcontext "todoshnik/internal/auth/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleAuth(ctx context.Context, user *tgbotapi.User) context.Context {
	cacheKey := tokenCacheKey(user.ID)

	// попытка забрать из кеша
	token, err := h.cache.Get(ctx, cacheKey).Result()
	if err == nil && token != "" {
		return authcontext.SetToken(ctx, token)
	}

	if token != "" {
		ctx = authcontext.SetToken(ctx, token)
		return ctx
	}

	// попытка сгенерировать через апи
	tokenInfo, err := h.authHandler.GetToken(ctx, authbot.TgLoginRequestDto{
		TgUserID: user.ID,
		Name:     user.UserName,
	})

	if err != nil {
		return ctx
	}

	// вставить в кеш
	if err := h.cache.Set(
		ctx,
		cacheKey,
		tokenInfo.AccessToken,
		botTokenTtl,
	).Err(); err != nil {
		h.logger.Printf("cache set error: %v", err)
	}

	// вставить в контекст
	ctx = authcontext.SetToken(ctx, tokenInfo.AccessToken)

	return ctx
}

func (h *Handler) handleWelcome(ctx context.Context, user *tgbotapi.User) error {
	cacheKey := tokenCacheKey(user.ID)
	tokenInfo, err := h.authHandler.SignInUser(ctx, authbot.TgLoginRequestDto{
		TgUserID: user.ID,
		Name:     user.UserName,
	})
	if err != nil {
		return err
	}

	if err := h.cache.Set(
		ctx,
		cacheKey,
		tokenInfo.AccessToken,
		botTokenTtl,
	).Err(); err != nil {
		h.logger.Printf("cache set error: %v", err)
	}

	return nil
}

func tokenCacheKey(userID int64) string {
	return fmt.Sprintf("user:%d:tg-api-token-key", userID)
}
