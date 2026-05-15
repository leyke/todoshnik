package session

import (
	"context"
	"todoshnik/internal/bot/tg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Storage) StartCommand(
	ctx context.Context,
	user *tgbotapi.User,
	command tg.Command,
) error {
	return s.Set(ctx, user.ID, command, tg.StateWait)
}

func (s *Storage) FinishCommand(
	ctx context.Context,
	user *tgbotapi.User,
	command tg.Command,
) error {
	return s.Set(ctx, user.ID, command, tg.StateComplete)
}
