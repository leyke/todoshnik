package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) send(c tgbotapi.Chattable) {
	if c == nil {
		return
	}

	_, err := h.bot.Send(c)
	if err != nil {
		h.logger.Println(err)
	}
}

func (h *Handler) editTgMessage(chatID int64, messageID int, newText string) {
	editMsg := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		newText,
	)
	h.bot.Send(editMsg)
}

func (h *Handler) deleteTgMessage(chatID int64, messageID int) {
	deleteMsg := tgbotapi.NewDeleteMessage(
		chatID,
		messageID,
	)
	h.bot.Send(deleteMsg)
}
