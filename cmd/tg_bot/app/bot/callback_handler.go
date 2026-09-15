package bot

import (
	"context"
	"encoding/json"

	"todoshnik/cmd/tg_bot/app/bot/command"
	"todoshnik/cmd/tg_bot/app/bot/response"

	boterrors "todoshnik/cmd/tg_bot/app/bot/errors"
	client "todoshnik/internal/infrastructure/api_client"
	authcontext "todoshnik/internal/infrastructure/context_manager/auth"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackQuery struct {
	Command command.Name      `json:"command"`
	Payload map[string]string `json:"payload"`
}

type CallbackHandler func(
	ctx context.Context,
	update tgbotapi.Update,
	callback *CallbackQuery,
) tgbotapi.Chattable

func (h *Handler) callbackHandlers() map[command.Name]CallbackHandler {
	return map[command.Name]CallbackHandler{
		command.CommandTaskDone:   h.callbackTaskDone,
		command.CommandTaskDelete: h.callbackTaskDelete,
	}
}

func (h *Handler) handleCallbackUpdate(ctx context.Context, update tgbotapi.Update) tgbotapi.Chattable {
	query := update.CallbackQuery
	callback, err := decodeJSON(update.CallbackQuery)
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}

	_, err = h.bot.Request(tgbotapi.NewCallback(query.ID, ""))
	if err != nil {
		h.logger.Println(err)
	}

	handlers := h.callbackHandlers()

	handler, ok := handlers[callback.Command]
	if !ok {
		return response.NewError(update.Message.Chat.ID, boterrors.ErrUnknownMethod)
	}

	// попытка авторизоваться
	token, err := h.handleAuth(ctx, query.From)
	if err != nil || token == "" {
		h.logger.Printf("ошибка авторизации: %v", err)
		return response.NewError(update.Message.Chat.ID, client.ErrUnAuth)
	}

	// вставить в контекст
	ctx = authcontext.SetToken(ctx, token)

	return handler(ctx, update, callback)
}

func (h *Handler) callbackTaskDone(
	ctx context.Context,
	update tgbotapi.Update,
	callback *CallbackQuery,
) tgbotapi.Chattable {
	query := update.CallbackQuery
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, "")

	taskId, ok := getTaskIDFromCallback(callback)
	if !ok {
		return response.NewError(update.Message.Chat.ID, boterrors.ErrInvalidTaskID)
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
	callback *CallbackQuery,
) tgbotapi.Chattable {
	query := update.CallbackQuery
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, "")

	taskId, ok := getTaskIDFromCallback(callback)
	if !ok {
		return response.NewError(update.Message.Chat.ID, boterrors.ErrInvalidTaskID)
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

func decodeJSON(callback *tgbotapi.CallbackQuery) (*CallbackQuery, error) {
	var query CallbackQuery

	err := json.Unmarshal([]byte(callback.Data), &query)
	if err != nil {
		return nil, err
	}
	return &query, nil
}

func getTaskIDFromCallback(callback *CallbackQuery) (string, bool) {
	taskID, ok := callback.Payload["task_id"]
	return taskID, ok
}
