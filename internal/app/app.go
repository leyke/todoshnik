package app

import (
	"errors"
	"fmt"
	"log"

	"todoshnik/internal/config"
	"todoshnik/internal/domains/task"
	"todoshnik/internal/domains/token"
	"todoshnik/internal/domains/user"
	"todoshnik/internal/infrastructure/db"

	taskrepo "todoshnik/internal/infrastructure/db/repository/task"
	tokenrepo "todoshnik/internal/infrastructure/db/repository/token"
	userrepo "todoshnik/internal/infrastructure/db/repository/user"

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

	return &App{
		Logger:   log,
		Cache:    redisBase,
		DB:       dataBase,
		Services: newServices(dataBase, cfg),
	}, nil
}

func newServices(dataBase *gorm.DB, cfg config.Config) *Services {
	// Репозитории
	taskRepo := taskrepo.NewDbRepository(dataBase)
	userRepo := userrepo.NewDbRepository(dataBase)
	tokenRepo := tokenrepo.NewDbRepository(dataBase)

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
