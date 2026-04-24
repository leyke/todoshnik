package main

import (
	"log"
	"os"
	"todoshnik/internal/app"
	"todoshnik/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var logPath string = "./tmp/tg.log"

func main() {
	container := app.InitApp(logPath)
	defer container.LogFile.Close()

	botapi, err := tgbotapi.NewBotAPI(os.Getenv("telegram_token"))
	if err != nil {
		log.Panic(err)
	}

	botapi.Debug = os.Getenv("telegram_debug") == "1"

	bh := bot.NewBotHandler(container, botapi, container.Logger)
	bh.Run()
}
