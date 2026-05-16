package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"todoshnik/internal/app"
	"todoshnik/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var logFileName string = "/tg.log"

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	container := app.InitApp(logFileName)

	botapi, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	botapi.Debug = os.Getenv("TELEGRAM_DEBUG") == "1"

	bh := bot.NewHandler(container, botapi)

	go func() {
		if err := bh.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	log.Println("Сигнал выкл получен")

	bh.Shutdown()

	container.Cache.Close()
	container.LogFile.Close()
}
