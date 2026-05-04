package main

import (
	"log"
	"os"
	"todoshnik/internal/app"
	"todoshnik/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var logFile string = "/tg.log"

func main() {
	container := app.InitApp(logFile)
	defer container.LogFile.Close()

	botapi, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	botapi.Debug = os.Getenv("TELEGRAM_DEBUG") == "1"

	bh := bot.NewBotHandler(container, botapi, container.Logger)
	bh.Run()
}
