package main

import (
	"log"
	"os"
	"todoshnik/internal/api"
	"todoshnik/internal/app"
)

func main() {
	app := app.InitApp()
	logFile, err := os.OpenFile("./tmp/api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	api := api.NewAPIHandler(app.TaskService, logFile)
	defer api.Close()

	api.Run()
}
