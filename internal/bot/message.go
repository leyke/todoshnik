package bot

import (
	"context"
	"fmt"

	boterrors "todoshnik/internal/bot/errors"
	"todoshnik/internal/bot/response"
	"todoshnik/internal/bot/tg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandler func(
	ctx context.Context,
	update tgbotapi.Update,
	text string,
) tgbotapi.Chattable

func (h *Handler) messageHandlers() map[tg.Command]MessageHandler {
	return map[tg.Command]MessageHandler{
		tg.CommandAdd: h.messageTaskAdd,
	}
}

func (h *Handler) handleMessageUpdate(ctx context.Context, update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tgUser := update.Message.From

	lastState, ok := h.sessionStorage.Get(ctx, tgUser.ID)
	if !ok {
		msg.Text = "Я забыл на чем мы остановились, повтори ввод команды"
		return &msg
	}

	if lastState.State != tg.StateWait {
		msg.Text = "Я уже все сделал, начни новую команду"
		return &msg
	}

	handlers := h.messageHandlers()

	handler, ok := handlers[lastState.Command]
	if !ok {
		h.logger.Println(boterrors.ErrUnknownMethod)
		return response.NewError(update.Message.Chat.ID, boterrors.ErrUnknownMethod)
	}

	// попытка авторизоваться
	ctx = h.handleAuth(ctx, tgUser)

	return handler(ctx, update, update.Message.Text)
}

func (h *Handler) messageTaskAdd(
	ctx context.Context,
	update tgbotapi.Update,
	text string,
) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tgUser := update.Message.From

	task, err := h.taskHandler.Add(ctx, text)
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}

	msg.Text = fmt.Sprintf("Добавил: %s", task.Title)
	h.sessionStorage.FinishCommand(ctx, tgUser, tg.CommandAdd)

	return msg
}
