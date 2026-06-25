package bot

import (
	"context"
	"fmt"

	"todoshnik/cmd/tg_bot/app/bot/command"
	"todoshnik/cmd/tg_bot/app/bot/response"

	boterrors "todoshnik/cmd/tg_bot/app/bot/errors"
	client "todoshnik/internal/infrastructure/api_client"
	authcontext "todoshnik/internal/infrastructure/context_manager/auth"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandler func(
	ctx context.Context,
	update tgbotapi.Update,
	text string,
) tgbotapi.Chattable

func (h *Handler) messageHandlers() map[command.Name]MessageHandler {
	return map[command.Name]MessageHandler{
		command.CommandAdd: h.messageTaskAdd,
	}
}

func (h *Handler) handleMessageUpdate(ctx context.Context, update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tgUser := update.Message.From

	lastCommand, ok := h.commandStorage.GetLastCommand(ctx, tgUser.ID)
	if !ok {
		msg.Text = "Я забыл на чем мы остановились, повтори ввод команды"
		return &msg
	}

	if lastCommand.State != command.StateWait {
		msg.Text = "Я уже все сделал, начни новую команду"
		return &msg
	}

	handlers := h.messageHandlers()

	handler, ok := handlers[lastCommand.Name]
	if !ok {
		h.logger.Println(boterrors.ErrUnknownMethod)
		return response.NewError(update.Message.Chat.ID, boterrors.ErrUnknownMethod)
	}

	// попытка авторизоваться
	token, err := h.handleAuth(ctx, tgUser)
	if err != nil || token == "" {
		h.logger.Printf("ошибка авторизации: %v", err)
		return response.NewError(update.Message.Chat.ID, client.ErrUnAuth)
	}

	// вставить в контекст
	ctx = authcontext.SetToken(ctx, token)

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
	err = h.commandStorage.FinishCommand(ctx, tgUser.ID, command.CommandAdd)
	if err != nil {
		h.logger.Println(err)
	}

	return msg
}
