package response

import (
	"encoding/json"
	"strconv"

	"todoshnik/cmd/tg_bot/app/bot/callback"
	"todoshnik/cmd/tg_bot/app/bot/command"
	"todoshnik/internal/domains/task"

	taskbot "todoshnik/cmd/tg_bot/app/task"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewTaskMsg(chatID int64, task *task.Task) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = taskbot.GetTaskRowText(*task)

	// добавим кнопки для управления
	var btns []InlineKeyboardBtn
	payload := map[string]string{
		"task_id": strconv.Itoa(task.ID),
	}

	callbackDone, err := json.Marshal(callback.CallbackQuery{
		Command: command.CommandTaskDone,
		Payload: payload,
	})
	if err != nil {
		return nil, err
	}

	callbackDelete, err := json.Marshal(callback.CallbackQuery{
		Command: command.CommandTaskDelete,
		Payload: payload,
	})
	if err != nil {
		return nil, err
	}

	// Кнопка снять/добавить галочку (выполнен или нет)
	btns = append(btns, InlineKeyboardBtn{
		Text:     taskbot.GetStatusButtonText(*task),
		Callback: string(callbackDone),
	})
	// Кнопка удалить таск
	btns = append(btns, InlineKeyboardBtn{
		Text:     taskbot.GetDeleteButtonText(),
		Callback: string(callbackDelete),
	})

	msg.ReplyMarkup = NewKeyboard(btns)

	return &msg, nil
}
