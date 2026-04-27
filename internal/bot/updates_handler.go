package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"todoshnik/internal/app"
	"todoshnik/internal/bot/task"
	"todoshnik/internal/bot/tg"
	"todoshnik/internal/bot/user"
	apperrors "todoshnik/internal/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	TaskHandler  *task.Handler
	UserHandler  *user.Handler
	StateStorage *user.StateStorage
	bot          *tgbotapi.BotAPI
	logger       *log.Logger
}

func NewBotHandler(container *app.App, bot *tgbotapi.BotAPI, logger *log.Logger) *BotHandler {
	return &BotHandler{
		TaskHandler:  task.NewHandler(container.TaskService, logger),
		UserHandler:  user.NewHandler(container.UserService),
		StateStorage: user.NewStateStorage(),
		bot:          bot,
		logger:       logger,
	}
}

func (bh *BotHandler) Run() {
	log.Printf("Authorized on account %s", bh.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bh.bot.GetUpdatesChan(u)

	for update := range updates {

		var msg *tgbotapi.MessageConfig

		if update.CallbackQuery != nil {
			msg = bh.handleCallback(update)
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				msg = bh.handleCommand(update)
			} else {
				msg = bh.handleMessage(update)
			}
		} else {
			continue
		}

		if _, err := bh.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func (bh *BotHandler) handleCallback(update tgbotapi.Update) *tgbotapi.MessageConfig {
	query := update.CallbackQuery
	tgUser := query.From
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, "")

	bh.bot.Request(tgbotapi.NewCallback(query.ID, ""))

	var callback *tg.CallbackQuery
	err := json.Unmarshal([]byte(query.Data), &callback)
	if err != nil {
		msg.Text = "Возникла непредвиденная ошибка"
		fmt.Println(query.Data)
		fmt.Println(err.Error())
		bh.logger.Println(err)
		return &msg
	}

	appUser, err := bh.UserHandler.GetAppUser(tgUser)
	if err != nil {
		msg.Text = "Я тебя забыл, давай познакомимся еще раз /restart"
		return &msg
	}

	switch callback.Command {
	case tg.СommandTaskDone:
		err := bh.TaskHandler.DoneTask(appUser.ID, callback.Payload["task_id"])
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				msg.Text = err.Error()
			} else {
				msg.Text = "Возникла непредвиденная ошибка"
				fmt.Println(err.Error())
				bh.logger.Println(err)
			}
		} else {
			msg.Text = "Статус обновлен"
		}
	default:
		msg.Text = "Я хз что это такое, если бы я знал что это такое, я бы помог"
	}

	return &msg
}

func (bh *BotHandler) handleCommand(update tgbotapi.Update) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tgUser := update.Message.From

	command := tg.Command(update.Message.Command())
	args := update.Message.CommandArguments()
	switch command {
	case tg.CommandStart, tg.CommandRestart:
		bh.UserHandler.AddUser(tgUser)
		msg.Text = "Привет, я готов запоминать задачи, начни с /add"
	case tg.CommandHelp:
		msg.Text = "Я могу /add, /list, /status и /restart."
	case tg.CommandStatus:
		msg.Text = "Я ОК."
	case tg.CommandAdd:
		appUser, err := bh.UserHandler.GetAppUser(tgUser)
		if err != nil {
			msg.Text = "Я тебя забыл, давай познакомимся еще раз /restart"
			return &msg
		}

		if args == "" {
			bh.startCommandHandling(tgUser, tg.CommandAdd)
			msg.Text = "Напиши задачу и я её запомню!"
		} else {
			bh.TaskHandler.AddTask(appUser.ID, args)
			msg.Text = fmt.Sprintf("Добавил: %v", args)
		}
	case tg.СommandTaskList:
		appUser, err := bh.UserHandler.GetAppUser(tgUser)
		if err != nil {
			msg.Text = "Я тебя забыл, давай познакомимся еще раз /restart"
			return &msg
		}

		count := bh.TaskHandler.SendTaskList(bh.bot, update.Message.Chat.ID, appUser.ID, args)
		if count == 0 {
			msg.Text = "У меня пока нет твоих задач. Давай добавим /add"
		} else {
			msg.Text = "Вот твои задачи"
		}
	default:
		msg.Text = "Я не знаю такой команды"
	}

	return &msg
}

func (bh *BotHandler) handleMessage(update tgbotapi.Update) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tgUser := update.Message.From

	lastState, ok := bh.StateStorage.Get(tgUser.ID)
	if !ok {
		msg.Text = "Я забыл на чем мы остановились, повтори ввод команды"
		return &msg
	}
	fmt.Println(lastState)
	if lastState.State != tg.StateWait {
		msg.Text = "Я уже все сделал, начни новую команду"
		return &msg
	}

	appUser, err := bh.UserHandler.GetAppUser(tgUser)
	if err != nil {
		msg.Text = "Я тебя забыл, давай познакомимся еще раз /restart"
		return &msg
	}

	switch lastState.Command {
	case tg.CommandAdd:
		task, err := bh.TaskHandler.AddTask(appUser.ID, update.Message.Text)
		if err != nil {
			if errors.Is(err, apperrors.ErrNotValidate) {
				msg.Text = fmt.Sprintf("Возникла ошибка: %v", err.Error())
			} else {
				msg.Text = "Возникла непредвиденная ошибка"
				fmt.Println(err.Error())
				bh.logger.Println(err)
			}
		} else {
			msg.Text = fmt.Sprintf("Добавил: %s", task.Title)
			bh.finishCommandHandling(tgUser, tg.CommandAdd)
		}
	default:
		msg.Text = "Для этой команды я уже ничего не могу сделать, начни новую"
	}

	return &msg
}

func (bh BotHandler) startCommandHandling(user *tgbotapi.User, command tg.Command) {
	bh.StateStorage.Set(user.ID, command, tg.StateWait)
}

func (bh BotHandler) finishCommandHandling(user *tgbotapi.User, command tg.Command) {
	bh.StateStorage.Set(user.ID, command, tg.StateComplete)
}
