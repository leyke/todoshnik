package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"todoshnik/internal/config"
	"todoshnik/internal/domains/task"
	"todoshnik/internal/domains/token"
	"todoshnik/internal/domains/user"

	"todoshnik/internal/infrastructure/db"
	"todoshnik/internal/infrastructure/db/transaction"

	taskrepo "todoshnik/internal/infrastructure/db/repository/task"
	tokenrepo "todoshnik/internal/infrastructure/db/repository/token"
	userrepo "todoshnik/internal/infrastructure/db/repository/user"

	rdb "todoshnik/internal/infrastructure/redis"

	"github.com/redis/go-redis/v9"
)

type App struct {
	DB         *sql.DB
	Cache      *redis.Client
	Logger     *log.Logger
	Transactor *transaction.Transactor

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

	dataBase, err := db.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	transactor := transaction.NewTransactor(dataBase)

	return &App{
		Logger:     log,
		Cache:      redisBase,
		DB:         dataBase,
		Transactor: transactor,
		Services:   newServices(dataBase, cfg),
	}, nil
}

func newServices(dataBase *sql.DB, cfg config.Config) *Services {
	// Репозитории
	taskRepo := taskrepo.NewRepository(dataBase)
	userRepo := userrepo.NewRepository(dataBase)
	tokenRepo := tokenrepo.NewRepository(dataBase)

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
		if err := app.DB.Close(); err != nil {
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
