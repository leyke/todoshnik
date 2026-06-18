package app

import (
	"errors"
	"fmt"
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

	Services *Services
}

type Services struct {
	TaskService  *task.Service
	UserService  *user.Service
	TokenService *token.Service
}

func InitApp(cfg config.Config) (*App, error) {
	log := NewLogger()

	redisBase, err := rdb.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("init redis: %w", err)
	}

	dataBase, err := db.NewGormDb(cfg)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	// TODO тут определиться со структурой
	// 1. За то что есть структура App- однозначно лайк, но я бы группировал более гранулярно по уровням, потому что
	// уровни приложения не должны пересекаться и это проявляется на уровне инициализации: инфраструктурные зависимости:
	// logger, logFile это core зависимости, redis зависит от logger по идее и это инфраструктурная зависимость.
	// TaskService, UserService, TokenService зависят от логера и от инфраструктурной зависимости (бд, редиса).
	return &App{
		Services: newServices(dataBase, cfg),
		Logger:   log,
		Cache:    redisBase,
		DB:       dataBase,
	}, nil
}

func newServices(dataBase *gorm.DB, cfg config.Config) *Services {
	// Репозитории
	taskRepo := task.NewDbRepository(dataBase)
	userRepo := user.NewDbRepository(dataBase)
	tokenRepo := token.NewDbRepository(dataBase)

	return &Services{
		TaskService: task.NewService(taskRepo),
		UserService: user.NewService(userRepo),
		TokenService: token.NewService(tokenRepo, token.Config{
			Secret: cfg.App.TokenSecret,
			Ttl:    cfg.App.TokenTtl,
		}),
	}
}

func (app *App) Close() error {
	var errs []error

	if app.DB != nil {
		sqlDB, err := app.DB.DB()
		if err != nil {
			errs = append(errs, fmt.Errorf("get sql db: %w", err))
		} else if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close sql db: %w", err))
		}
	}

	if app.Cache != nil {
		if err := app.Cache.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close redis: %w", err))
		}
	}

	return errors.Join(errs...)
}
