package main

import (
	"todoshnik/internal/api"
	"todoshnik/internal/app"
)

func main() {
	app := app.InitApp()

	api := api.NewAPIHandler(app.TaskService, "./tmp/api.log")
	api.Run()
}
