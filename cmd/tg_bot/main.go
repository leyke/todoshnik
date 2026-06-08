package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"todoshnik/internal/app"
	"todoshnik/internal/bot"
	"todoshnik/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	container, err := app.InitApp(cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	botapi, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatalf("create bot API: %v", err)
	}
	botapi.Debug = cfg.Telegram.Debug

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
}
