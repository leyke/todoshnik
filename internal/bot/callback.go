package bot

import (
	"context"

	"todoshnik/internal/bot/response"
	"todoshnik/internal/bot/tg"
	"todoshnik/internal/task"

	boterrors "todoshnik/internal/bot/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler func(
	ctx context.Context,
	update tgbotapi.Update,
	callback *tg.CallbackQuery,
) tgbotapi.Chattable

func (h *Handler) callbackHandlers() map[tg.Command]CallbackHandler {
	return map[tg.Command]CallbackHandler{
		tg.CommandTaskDone:   h.callbackTaskDone,
		tg.CommandTaskDelete: h.callbackTaskDelete,
	}
}

func (h *Handler) handleCallbackUpdate(ctx context.Context, update tgbotapi.Update) tgbotapi.Chattable {
	query := update.CallbackQuery
	callback, err := tg.DecodeJSON(update.CallbackQuery)
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}
	h.bot.Request(tgbotapi.NewCallback(query.ID, ""))

	handlers := h.callbackHandlers()

	handler, ok := handlers[callback.Command]
	if !ok {
		return response.NewError(update.Message.Chat.ID, boterrors.ErrUnknownMethod)
	}

	// попытка авторизоваться
	ctx = h.handleAuth(ctx, query.From)

	return handler(ctx, update, callback)
}

func (h *Handler) callbackTaskDone(
	ctx context.Context,
	update tgbotapi.Update,
	callback *tg.CallbackQuery,
) tgbotapi.Chattable {
	query := update.CallbackQuery
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, "")

	taskId, ok := tg.GetTaskID(callback)
	if !ok {
		return response.NewError(update.Message.Chat.ID, task.ErrInvalidTaskID)
	}

	taskRowText, err := h.taskHandler.DoneTask(ctx, taskId)
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)

	}
	if taskRowText != "" {
		h.editTgMessage(query.Message.Chat.ID, query.Message.MessageID, taskRowText)
	}

	msg.Text = "Статус обновлен"
	return &msg
}

func (h *Handler) callbackTaskDelete(
	ctx context.Context,
	update tgbotapi.Update,
	callback *tg.CallbackQuery,
) tgbotapi.Chattable {
	query := update.CallbackQuery
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, "")

	taskId, ok := tg.GetTaskID(callback)
	if !ok {
		return response.NewError(update.Message.Chat.ID, task.ErrInvalidTaskID)
	}

	err := h.taskHandler.DeleteTask(ctx, taskId)
	if err != nil {
		h.logger.Println(err)
		return response.NewError(query.Message.Chat.ID, err)
	}

	msg.Text = "Задача удалена"
	h.deleteTgMessage(query.Message.Chat.ID, query.Message.MessageID)
	return &msg
}
