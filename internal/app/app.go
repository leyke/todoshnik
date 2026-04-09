package app

import (
	"todoshnik/internal/service"
	"todoshnik/internal/storage"
)

type App struct {
	TaskService *service.TaskService
}

func InitApp() *App {
	storage := storage.NewFileStorage("./tmp/tasks.json")
	tm, err := service.NewTaskService(storage)
	if err != nil {
		panic(err)
	}
	return &App{TaskService: tm}
}
