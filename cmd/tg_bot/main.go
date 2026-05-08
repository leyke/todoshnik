package main

import (
	"log"
	"os"
	"todoshnik/internal/app"
	"todoshnik/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var logFileName string = "/tg.log"

func main() {
	container := app.InitApp(logFileName)
	defer container.LogFile.Close()
	defer container.Cache.Close()

	botapi, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	botapi.Debug = os.Getenv("TELEGRAM_DEBUG") == "1"

	bh := bot.NewBotHandler(container, botapi)
	bh.Run()
}
