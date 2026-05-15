package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) dispatchUpdate(update tgbotapi.Update) {
	ctx := context.Background()

	var msg tgbotapi.Chattable

	switch {
	case update.CallbackQuery != nil:
		msg = h.handleCallbackUpdate(ctx, update)

	case update.Message != nil && update.Message.IsCommand():
		msg = h.handleCommand(ctx, update)

	case update.Message != nil:
		msg = h.handleMessageUpdate(ctx, update)
	}

	if msg != nil {
		h.send(msg)
	}
}
