package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"todoshnik/internal/api"
	"todoshnik/internal/app"
	"todoshnik/internal/config"
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

	apiHandler := api.NewAPIHandler(container)

	go func() {
		if err := apiHandler.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := apiHandler.Shutdown(shutdownCtx); err != nil {
		log.Println(err)
	}

	if err := container.Cache.Close(); err != nil {
		log.Println(err)
	}
}
