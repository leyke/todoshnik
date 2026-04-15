package main

import (
	"todoshnik/internal/api"
	"todoshnik/internal/app"
)

var logPath string = "./tmp/api.log"

func main() {
	container := app.InitApp()
	log, logFile := app.NewLogger(logPath)

	api := api.NewAPIHandler(container, log)
	defer logFile.Close()

	api.Run()
}
