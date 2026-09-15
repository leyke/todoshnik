package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"todoshnik/cmd/tg_bot/app"
	"todoshnik/cmd/tg_bot/app/bot"
	"todoshnik/internal/config"
)

var (
	AppVersion = "dev"
	CommitHash = "unknown"
)

func main() {
	fmt.Println("Version:", AppVersion)
	fmt.Println("Commit :", CommitHash)

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

	bh := bot.NewHandler(
		container,
		bot.Config{
			ApiURL:            cfg.Telegram.BotApiUrl,
			BotServiceToken:   cfg.Telegram.ServiceToken,
			SemaphoreSize:     cfg.Telegram.SemaphoreSize,
			ApiRequestTimeout: cfg.Telegram.ApiRequestTimeOut,
		},
	)

	go func() {
		if err := bh.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	log.Println("Сигнал выкл получен")

	bh.Shutdown()

	if err := container.Close(); err != nil {
		log.Println(err)
	}
}
