package response

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type InlineKeyboardBtn struct {
	Text     string
	Callback string
}

// Создает одно строковую клавиатуру
func NewKeyboard(btns []InlineKeyboardBtn) tgbotapi.InlineKeyboardMarkup {
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(btns))
	for _, item := range btns {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(item.Text, item.Callback))
	}

	// обертка в строки, пока не требуется больше одной
	rows := tgbotapi.NewInlineKeyboardRow(row...)

	return tgbotapi.NewInlineKeyboardMarkup(rows)
}
