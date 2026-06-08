package main

import (
	"log"
	"os"

	"todoshnik/internal/app"
	"todoshnik/internal/cli"
	"todoshnik/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	container, err := app.InitApp(cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	cli := cli.NewHandler(container.TaskService, container.TokenService)
	cli.Run(os.Args)
}
