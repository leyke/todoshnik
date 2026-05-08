package task

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"todoshnik/internal/bot/response"
	"todoshnik/internal/bot/tg"
	"todoshnik/internal/client"
	"todoshnik/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	api    *client.ApiClient
	logger *log.Logger
}

func NewHandler(api *client.ApiClient, l *log.Logger) *Handler {
	return &Handler{api: api, logger: l}
}

func (h Handler) AddTask(ctx context.Context, taskTitle string) (*domain.Task, error) {
	response, err := h.api.Post(ctx, "/tasks", CreateTaskRequest{
		Title: taskTitle,
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var task *domain.Task

	err = json.NewDecoder(response.Body).Decode(&task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (h Handler) SendTaskList(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, status string) (int, error) {
	params := url.Values{}
	if status != "" {
		params.Add("status", status)
	}

	response, err := h.api.Get(ctx, "/tasks", params)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var tasks []*domain.Task

	err = json.NewDecoder(response.Body).Decode(&tasks)
	if err != nil {
		return 0, err
	}

	if len(tasks) == 0 {
		return 0, nil
	}

	messageCount := 0
	for _, task := range tasks {
		if h.sendTask(bot, chatID, task) {
			messageCount++
		}
	}

	return messageCount, nil
}

func (h Handler) DoneTask(ctx context.Context, taskID string) (string, error) {
	response, err := h.api.Post(ctx, "/tasks/"+taskID+"/done", nil)
	if err != nil {
		return "", err
	}

	var task domain.Task

	err = json.NewDecoder(response.Body).Decode(&task)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return getTaskRowText(task), nil
}

func (h Handler) DeleteTask(ctx context.Context, taskID string) error {
	response, err := h.api.Delete(ctx, "/tasks/"+taskID+"")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func (h Handler) sendTask(bot *tgbotapi.BotAPI, chatID int64, task *domain.Task) bool {
	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = getTaskRowText(*task)

	// добавим кнопки для управления
	var btns []response.InlineKeyboardBtn
	payload := map[string]string{
		"task_id": strconv.Itoa(task.ID),
	}

	сallbackDone, err := json.Marshal(tg.CallbackQuery{
		Command: tg.СommandTaskDone,
		Payload: payload,
	})
	сallbackDelete, err := json.Marshal(tg.CallbackQuery{
		Command: tg.CommandTaskDelete,
		Payload: payload,
	})

	if err != nil {
		h.logger.Println("sendTask | Ошибка кодирования payloadData", err)
		return false
	}

	// Кнопка снять/добавить галочку (выполнен или нет)
	btns = append(btns, response.InlineKeyboardBtn{
		Text:     getStatusButtonText(*task),
		Callback: string(сallbackDone),
	})
	// Кнопка удалить таск
	btns = append(btns, response.InlineKeyboardBtn{
		Text:     getDeleteButtonText(),
		Callback: string(сallbackDelete),
	})
	msg.ReplyMarkup = response.NewKeyboard(btns)

	if _, err := bot.Send(msg); err != nil {
		h.logger.Println("sendTask | Ошибка отправки сообщения", err)
		return false
	}

	return true
}
