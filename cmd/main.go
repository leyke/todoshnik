package main

import (
	"fmt"
	"os"

	"todoshnik/internal/cli"
	"todoshnik/internal/task"
)

func main() {
	os.MkdirAll("tmp", 0755)

	storage := task.NewFileStorage("./tmp/tasks.json")
	tm, err := task.NewTaskManager(storage)
	if err != nil {
		fmt.Printf("Ошибка инициализации менеджера задач: %v\n", err)
		return
	}

	cli.Run(tm)
}
