package main

import (
	"os"
	"todoshnik/internal/app"
	"todoshnik/internal/cli"
)

var logFileName string = "/cli.log"

func main() {
	container := app.InitApp(logFileName)
	defer container.LogFile.Close()

	cli := cli.NewHandler(container.TaskService, container.TokenService)
	cli.Run(os.Args)
}
