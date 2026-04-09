package main

import (
	"todoshnik/internal/app"
	"todoshnik/internal/cli"
)

func main() {
	appIntance := app.InitApp()
	cli := cli.NewCLIHandler(appIntance.TaskService)
	cli.Run()
}
