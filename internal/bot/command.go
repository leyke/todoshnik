package bot

import (
	"context"
	"fmt"
	"todoshnik/internal/bot/response"
	"todoshnik/internal/bot/tg"
	apperrors "todoshnik/internal/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler func(
	ctx context.Context,
	update tgbotapi.Update,
	args string,
) tgbotapi.Chattable

func (h *Handler) commandHandlers() map[tg.Command]CommandHandler {
	return map[tg.Command]CommandHandler{
		tg.CommandStart:    h.cmdStart,
		tg.CommandRestart:  h.cmdStart,
		tg.CommandHelp:     h.cmdHelp,
		tg.CommandStatus:   h.cmdStatus,
		tg.CommandAdd:      h.cmdAdd,
		tg.CommandTaskList: h.cmdTaskList,
	}
}

func (h *Handler) handleCommand(
	ctx context.Context,
	update tgbotapi.Update,
) tgbotapi.Chattable {
	command := tg.Command(update.Message.Command())

	args := update.Message.CommandArguments()

	handlers := h.commandHandlers()

	handler, ok := handlers[command]
	if !ok {
		h.logger.Println(apperrors.ErrUnknownMethod)
		return response.NewError(update.Message.Chat.ID, apperrors.ErrUnknownMethod)
	}

	// авторизуем пользователя для всех команд, кроме стартовой
	if command != tg.CommandStart && command != tg.CommandRestart {
		ctx = h.handleAuth(ctx, update.Message.From)
	}

	return handler(ctx, update, args)
}

func (h *Handler) cmdStart(
	ctx context.Context,
	update tgbotapi.Update,
	args string,
) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tgUser := update.Message.From
	err := h.handleWelcome(ctx, tgUser)
	if err != nil {
		h.logger.Println(apperrors.ErrUnknownMethod)
		return response.NewError(update.Message.Chat.ID, err)
	}

	msg.Text = "Привет, я готов запоминать задачи, начни с /add"
	return &msg
}

func (h *Handler) cmdAdd(
	ctx context.Context,
	update tgbotapi.Update,
	args string,
) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tgUser := update.Message.From

	if args == "" {
		h.sessionStorage.StartCommand(ctx, tgUser, tg.CommandAdd)
		msg.Text = "Напиши задачу и я её запомню!"
		return &msg
	}

	task, err := h.taskHandler.Add(ctx, args)
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}

	msg.Text = fmt.Sprintf("Добавил: %s", task.Title)

	return &msg
}

func (h *Handler) cmdHelp(
	ctx context.Context,
	update tgbotapi.Update,
	args string,
) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Я могу /add, /tasklist, /status и /restart."

	return &msg
}

func (h *Handler) cmdStatus(
	ctx context.Context,
	update tgbotapi.Update,
	args string,
) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Я OK."
	return &msg
}

func (h *Handler) cmdTaskList(
	ctx context.Context,
	update tgbotapi.Update,
	args string,
) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tasks, err := h.taskHandler.List(ctx, args)

	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}

	if len(tasks) == 0 {
		msg.Text = "У меня пока нет твоих задач. Давай добавим /add"

		return &msg
	}

	for _, task := range tasks {
		taskMsg, err := response.NewTaskMsg(update.Message.Chat.ID, task)
		if err != nil {
			h.logger.Println("handleCommand | Ошибка создания сообщения с задачей", err)
			continue
		}
		h.bot.Send(taskMsg)
	}

	msg.Text = "Вот твои задачи"

	return &msg
}
