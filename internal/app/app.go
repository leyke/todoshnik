package app

import (
	"log"

	"todoshnik/internal/auth/token"
	"todoshnik/internal/config"
	"todoshnik/internal/infrastructure/db"
	"todoshnik/internal/task"
	"todoshnik/internal/user"

	rdb "todoshnik/internal/infrastructure/redis"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Cache  *redis.Client
	Logger *log.Logger

	TaskService  *task.Service
	UserService  *user.Service
	TokenService *token.Service
}

// TODO изучить работу с ошибками и wrapping ошибок !!!
func InitApp(cfg *config.Config) (*App, error) {
	log := NewLogger()

	dataBase, err := db.NewGormDb(cfg)
	if err != nil {
		return nil, err
	}

	redisBase, err := rdb.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Репозитории
	taskRepo := task.NewDbRepository(dataBase)
	userRepo := user.NewDbRepository(dataBase)
	tokenRepo := token.NewDbRepository(dataBase)

	// 1. За то что есть структура App- однозначно лайк, но я бы группировал более гранулярно по уровням, потому что
	// уровни приложения не должны пересекаться и это проявляется на уровне инициализации: инфраструктурные зависимости:
	// logger, logFile это core зависимости, redis зависит от logger по идее и это инфраструктурная зависимость.
	// TaskService, UserService, TokenService зависят от логера и от инфраструктурной зависимости (бд, редиса). Кстати
	// не вижу в Апп БД которую мы инициализировали. Все что мы инициализируем в Апп должно быть в Апп иначе это
	// утечка ресурса.
	// 2. Как писал выше должен быть какой-то app.Close() который будет завершать работу и заканчивать работу с
	// зависимостями и освобождать ресурсы
	return &App{
		TaskService:  task.NewService(taskRepo),
		UserService:  user.NewService(userRepo),
		TokenService: token.NewService(tokenRepo),
		Logger:       log,
		Cache:        redisBase,
		DB:           dataBase,
	}, nil
}
