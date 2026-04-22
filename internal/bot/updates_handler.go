package bot

import (
	"fmt"
	"log"
	"todoshnik/internal/app"
	"todoshnik/internal/bot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	TaskHandler handlers.TaskHandler
	UserHandler handlers.UserHandler
	bot         *tgbotapi.BotAPI
	logger      *log.Logger
}

func NewBotHandler(container *app.App, bot *tgbotapi.BotAPI, logger *log.Logger) *BotHandler {
	return &BotHandler{
		TaskHandler: *handlers.NewTaskHandler(container.TaskService),
		UserHandler: *handlers.NewUserHandler(container.UserService),
		bot:         bot,
		logger:      logger,
	}
}

func (bh *BotHandler) Run() {
	log.Printf("Authorized on account %s", bh.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bh.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			fmt.Printf("Получено сообщение от %s: %s\n", update.Message.From.UserName, update.Message.Text)
			fmt.Println(update.Message)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bh.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
