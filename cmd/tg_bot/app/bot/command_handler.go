package bot

import (
	"context"

	"todoshnik/cmd/tg_bot/app/bot/command"
	"todoshnik/cmd/tg_bot/app/bot/response"

	boterrors "todoshnik/cmd/tg_bot/app/bot/errors"
	client "todoshnik/internal/infrastructure/api_client"
	authcontext "todoshnik/internal/infrastructure/context_manager/auth"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler func(
	ctx context.Context,
	update tgbotapi.Update,
	args string,
) tgbotapi.Chattable

func (h *Handler) commandHandlers() map[command.Name]CommandHandler {
	return map[command.Name]CommandHandler{
		command.CommandStart:    h.cmdStart,
		command.CommandRestart:  h.cmdStart,
		command.CommandHelp:     h.cmdHelp,
		command.CommandStatus:   h.cmdStatus,
		command.CommandAdd:      h.cmdAdd,
		command.CommandTaskList: h.cmdTaskList,
	}
}

func (h *Handler) handleCommand(
	ctx context.Context,
	update tgbotapi.Update,
) tgbotapi.Chattable {
	c := command.Name(update.Message.Command())

	args := update.Message.CommandArguments()

	handlers := h.commandHandlers()

	handler, ok := handlers[c]
	if !ok {
		h.logger.Println(boterrors.ErrUnknownMethod)
		return response.NewError(update.Message.Chat.ID, boterrors.ErrUnknownMethod)
	}

	// авторизуем пользователя для всех команд, кроме стартовой
	if c != command.CommandStart && c != command.CommandRestart {
		token, err := h.handleAuth(ctx, update.Message.From)
		if err != nil {
			h.logger.Println(err)
			return response.NewError(update.Message.Chat.ID, client.ErrUnAuth)
		}

		// вставить в контекст
		ctx = authcontext.SetToken(ctx, token)
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
		h.logger.Printf("ошибка авторизации: %v", err)
		return response.NewError(update.Message.Chat.ID, client.ErrUnAuth)
	}

	text, err := response.RenderStart()
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}

	msg.Text = text

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
		err := h.commandStorage.StartCommand(ctx, tgUser.ID, command.CommandAdd)
		if err != nil {
			h.logger.Println(err)
			return response.NewError(update.Message.Chat.ID, err)
		}

		text, err := response.RenderAddRequest()
		if err != nil {
			h.logger.Println(err)
			return response.NewError(update.Message.Chat.ID, err)
		}

		msg.Text = text

		return &msg
	}

	task, err := h.taskHandler.Add(ctx, args)
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}

	text, err := response.RenderAddSuccess(task.Title)
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}
	msg.Text = text

	return &msg
}

func (h *Handler) cmdHelp(
	ctx context.Context,
	update tgbotapi.Update,
	args string,
) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	text, err := response.RenderHelp()
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}
	msg.Text = text

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
		text, err := response.RenderListEmpty()
		if err != nil {
			h.logger.Println(err)
			return response.NewError(update.Message.Chat.ID, err)
		}
		msg.Text = text

		return &msg
	}

	for _, task := range tasks {
		taskMsg, err := response.NewTaskMsg(update.Message.Chat.ID, task)
		if err != nil {
			h.logger.Println("handleCommand | Ошибка создания сообщения с задачей: %v", err)
			continue
		}
		_, err = h.bot.Send(taskMsg)
		if err != nil {
			h.logger.Printf("handleCommand | Ошибка отправки сообщения: %v", err)
			continue
		}
	}

	text, err := response.RenderListSuccess()
	if err != nil {
		h.logger.Println(err)
		return response.NewError(update.Message.Chat.ID, err)
	}
	msg.Text = text

	return &msg
}
