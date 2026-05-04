package app

import (
	"log"
	"os"
	"todoshnik/internal/service"

	"github.com/joho/godotenv"
)

type App struct {
	TaskService  *service.TaskService
	UserService  *service.UserService
	TokenService *service.AccessTokenService
	Logger       *log.Logger
	LogFile      *os.File
}

func InitApp(logFileName string) *App {
	_ = godotenv.Load()

	tmpDir := os.Getenv("TMP_DIR")
	os.MkdirAll(tmpDir, 0755)

	ts, err := service.NewTaskService()
	if err != nil {
		panic(err)
	}
	us, err := service.NewUserService()
	if err != nil {
		panic(err)
	}

	ats, err := service.NewAccessTokenService()
	if err != nil {
		panic(err)
	}

	log, logFile := NewLogger(tmpDir + logFileName)

	return &App{
		TaskService:  ts,
		UserService:  us,
		TokenService: ats,
		Logger:       log,
		LogFile:      logFile,
	}
}
