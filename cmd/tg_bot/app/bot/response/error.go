package response

import (
	"errors"

	boterrors "todoshnik/cmd/tg_bot/app/bot/errors"
	client "todoshnik/internal/infrastructure/api_client"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewError(chatID int64, err error) tgbotapi.Chattable {
	var msg tgbotapi.Chattable
	switch {
	case errors.Is(err, client.ErrNotFound):
		msg = tgbotapi.NewMessage(chatID, err.Error())
	case errors.Is(err, client.ErrUnAuth):
		msg = tgbotapi.NewMessage(chatID, "Я тебя забыл, давай познакомимся еще раз /restart")
	case errors.Is(err, boterrors.ErrUnknownMethod):
		msg = tgbotapi.NewMessage(chatID, "Я хз что это такое, если бы я знал что это такое, я бы помог /help")
	case errors.Is(err, boterrors.ErrInvalidTaskID):
		msg = tgbotapi.NewMessage(chatID, "Неверный ID задачи")
	default:
		msg = tgbotapi.NewMessage(chatID, "Возникла непредвиденная ошибка")
	}

	return msg
}
