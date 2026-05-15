package app

import (
	"log"
	"os"
	"todoshnik/internal/auth/token"
	"todoshnik/internal/infrastructure/db"
	rdb "todoshnik/internal/infrastructure/redis"
	"todoshnik/internal/task"
	"todoshnik/internal/user"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type App struct {
	TaskService  *task.Service
	UserService  *user.Service
	TokenService *token.Service
	Logger       *log.Logger
	LogFile      *os.File
	Cache        *redis.Client
}

func InitApp(logFileName string) *App {
	_ = godotenv.Load()

	tmpDir := os.Getenv("TMP_DIR")
	os.MkdirAll(tmpDir, 0755)

	log, logFile := NewLogger(tmpDir + logFileName)

	dataBase, err := db.NewGormDb()
	if err != nil {
		log.Fatal(err)
	}

	redisBase := rdb.NewClient()

	// Репозитории
	taskRepo := task.NewDbRepository(dataBase)
	userRepo := user.NewDbRepository(dataBase)
	tokenRepo := token.NewDbRepository(dataBase)

	return &App{
		TaskService:  task.NewService(taskRepo),
		UserService:  user.NewService(userRepo),
		TokenService: token.NewService(tokenRepo),
		Logger:       log,
		LogFile:      logFile,
		Cache:        redisBase,
	}
}
