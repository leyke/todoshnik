package main

import (
	"todoshnik/internal/api"
	"todoshnik/internal/app"
)

var logFile string = "/api.log"

func main() {
	container := app.InitApp(logFile)
	defer container.LogFile.Close()

	api := api.NewAPIHandler(container)
	api.Run()
}
