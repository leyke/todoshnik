package app

import (
	"log"
	"os"
	"todoshnik/internal/service"

	"github.com/joho/godotenv"
)

type App struct {
	TaskService *service.TaskService
	UserService *service.UserService
	Logger      *log.Logger
	LogFile     *os.File
}

func InitApp(logPath string) *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	os.MkdirAll(os.Getenv("tmp_dir"), 0755)

	ts, err := service.NewTaskService()
	if err != nil {
		panic(err)
	}
	us, err := service.NewUserService()
	if err != nil {
		panic(err)
	}
	log, logFile := NewLogger(logPath)

	return &App{
		TaskService: ts,
		UserService: us,
		Logger:      log,
		LogFile:     logFile,
	}
}
