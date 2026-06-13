package tg

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackQuery struct {
	Command Command           `json:"command"`
	Payload map[string]string `json:"payload"`
}

func GetTaskID(callback *CallbackQuery) (string, bool) {
	taskID, ok := callback.Payload["task_id"]
	return taskID, ok
}

func DecodeJSON(callback *tgbotapi.CallbackQuery) (*CallbackQuery, error) {
	var query CallbackQuery

	err := json.Unmarshal([]byte(callback.Data), &query)
	if err != nil {
		return nil, err
	}
	return &query, nil
}
