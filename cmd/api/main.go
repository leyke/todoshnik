package main

import (
	"todoshnik/internal/api"
	"todoshnik/internal/app"
)

var logFileName string = "/api.log"

func main() {
	container := app.InitApp(logFileName)
	defer container.LogFile.Close()

	api := api.NewAPIHandler(container)
	api.Run()
}
