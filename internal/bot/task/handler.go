package task

import (
	"encoding/json"
	"log"
	"strconv"
	"todoshnik/internal/bot/response"
	"todoshnik/internal/bot/tg"
	"todoshnik/internal/domain"
	"todoshnik/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	service *service.TaskService
	logger  *log.Logger
}

func NewHandler(s *service.TaskService, l *log.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

func (h Handler) AddTask(userID int, taskTitle string) (*domain.Task, error) {
	task, err := h.service.AddTask(taskTitle, &userID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (h Handler) SendTaskList(bot *tgbotapi.BotAPI, chatID int64, userID int, method string) int {
	tasks := h.service.ListTasks(method, &userID)
	if len(tasks) == 0 {
		return 0
	}

	messageCount := 0
	for _, task := range tasks {
		if h.sendTask(bot, chatID, task) {
			messageCount++
		}
	}

	return messageCount
}

func (h Handler) DoneTask(userID int, id string) (*domain.Task, error) {
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	task, err := h.service.MarkDone(taskID, &userID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (h Handler) sendTask(bot *tgbotapi.BotAPI, chatID int64, task *domain.Task) bool {
	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = getTaskRowText(*task)

	// добавим кнопки для управления
	var btns []response.InlineKeyboardBtn
	payload := map[string]string{
		"task_id": strconv.Itoa(task.ID),
	}
	сallback, err := json.Marshal(tg.CallbackQuery{
		Command: tg.СommandTaskDone,
		Payload: payload,
	})
	if err != nil {
		h.logger.Println("sendTask | Ошибка кодирования payloadData", err)
		return false
	}
	btns = append(btns, response.InlineKeyboardBtn{
		Text:     getStatusButtonText(*task),
		Callback: string(сallback),
	})
	msg.ReplyMarkup = response.NewKeyboard(btns)

	if _, err := bot.Send(msg); err != nil {
		h.logger.Println("sendTask | Ошибка отправки сообщения", err)
		return false
	}

	return true
}
