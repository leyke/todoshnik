package bot

import (
	"context"
	"fmt"

	"todoshnik/internal/bot/response"
	"todoshnik/internal/bot/tg"

	boterrors "todoshnik/internal/bot/errors"

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
		h.logger.Println(boterrors.ErrUnknownMethod)
		return response.NewError(update.Message.Chat.ID, boterrors.ErrUnknownMethod)
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
		h.logger.Println(boterrors.ErrUnknownMethod)
		return response.NewError(update.Message.Chat.ID, boterrors.ErrUnknownMethod)
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
	msg.Text = "Я могу /add, /list, /status и /restart."

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

	// мне бы не понравились куски тескта прямо в коде, надо резделять доменную логику, слой данных и представление
	// я бы для каждой команды/реплики сделал небольшой шаблон, организовал бы натягивание этих шаблонов на модель.
	// Да это более хлопотно, но в поддержке было бы гораздо проще (возможно). По крайней мере потом по коду искать
	// и пытаться понять почему то или иное сообщение вывелось будет довольно хлопотно. Изучи стандрантую шаблонизацию
	// в ГО: есть пакет https://pkg.go.dev/text/template и аналогичный html/template.
	// Надо проработать с моделями comand и под каждую команду завести свою структуду данных, необходимую для рендера
	// шаблона, шаблоны можно положить в подчиненный пакет
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
