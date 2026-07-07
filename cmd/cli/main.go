package main

import (
	"fmt"
	"log"
	"os"

	"todoshnik/cmd/tg_bot/app"
	"todoshnik/internal/config"

	cli "todoshnik/cmd/cli/app"
)

var (
	AppVersion = "dev"
	CommitHash = "unknown"
)

func main() {
	fmt.Println("Version:", AppVersion)
	fmt.Println("Commit :", CommitHash)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	container, err := app.InitApp(cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	cli := cli.NewHandler(container.Services)
	cli.Run(os.Args)

	if err := container.Close(); err != nil {
		log.Println(err)
	}
}
