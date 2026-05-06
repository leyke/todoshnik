package main

import (
	"todoshnik/internal/app"
	"todoshnik/internal/cli"
)

var logFileName string = "/cli.log"

func main() {
	container := app.InitApp(logFileName)
	defer container.LogFile.Close()

	cli := cli.NewCLIHandler(container.TaskService, container.TokenService)
	cli.Run()
}
