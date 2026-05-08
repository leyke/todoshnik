package app

import (
	"log"
	"os"
	"todoshnik/internal/infrastructure/db"
	rdb "todoshnik/internal/redis"
	tokenrepo "todoshnik/internal/repository/access_token"
	taskrepo "todoshnik/internal/repository/task"
	userrepo "todoshnik/internal/repository/user"
	"todoshnik/internal/service"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type App struct {
	TaskService  *service.TaskService
	UserService  *service.UserService
	TokenService *service.AccessTokenService
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
	taskRepo := taskrepo.NewTaskDbRepository(dataBase)
	userRepo := userrepo.NewUserDbRepository(dataBase)
	tokenRepo := tokenrepo.NewAccessTokenDbRepository(dataBase)

	return &App{
		TaskService:  service.NewTaskService(taskRepo),
		UserService:  service.NewUserService(userRepo),
		TokenService: service.NewAccessTokenService(tokenRepo),
		Logger:       log,
		LogFile:      logFile,
		Cache:        redisBase,
	}
}
