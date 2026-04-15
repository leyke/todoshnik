package main

import (
	"todoshnik/internal/app"
	"todoshnik/internal/cli"
)

func main() {
	container := app.InitApp()
	cli := cli.NewCLIHandler(container.TaskService)
	cli.Run()
}
