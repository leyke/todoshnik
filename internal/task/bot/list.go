package bot

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"todoshnik/internal/bot/response"
	"todoshnik/internal/bot/tg"
	"todoshnik/internal/task"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

	var tasks []*task.Task

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

func (h Handler) sendTask(bot *tgbotapi.BotAPI, chatID int64, task *task.Task) bool {
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
