package main

import (
	"todoshnik/internal/app"
	"todoshnik/internal/cli"
)

var logPath string = "./tmp/cli.log"

func main() {
	container := app.InitApp(logPath)
	defer container.LogFile.Close()

	cli := cli.NewCLIHandler(container.TaskService, container.TokenService)
	cli.Run()
}
