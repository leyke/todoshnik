package main

import (
	"todoshnik/internal/api"
	"todoshnik/internal/app"
)

var logPath string = "./tmp/api.log"

func main() {
	container := app.InitApp(logPath)
	defer container.LogFile.Close()

	api := api.NewAPIHandler(container)
	api.Run()
}
