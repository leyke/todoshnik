package main

import (
	"todoshnik/internal/app"
	"todoshnik/internal/cli"
)

var logFile string = "/cli.log"

func main() {
	container := app.InitApp(logFile)
	defer container.LogFile.Close()

	cli := cli.NewCLIHandler(container.TaskService, container.TokenService)
	cli.Run()
}
