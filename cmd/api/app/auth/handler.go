package auth

import (
	"context"
	"log"

	"todoshnik/internal/domains/token"
	"todoshnik/internal/domains/user"
)

type Handler struct {
	userService  UserService
	tokenService TokenService
	logger       *log.Logger
	transactor   Transactor
}

type TokenService interface {
	Get(ctx context.Context, rawToken string) (*token.Token, error)
	Add(ctx context.Context, user *user.User, device token.DeviceType) (string, error)
}

type UserService interface {
	Add(ctx context.Context, name string, login string, password string) (*user.User, error)
	AddFromTg(ctx context.Context, name string, telegramID int64) (*user.User, error)
	GetByLogin(ctx context.Context, login string) (*user.User, error)
	GetByTgId(ctx context.Context, userTgID int64) (*user.User, error)
	ValidatePassword(hash string, password string) (bool, error)
}

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}

func NewHandler(userService UserService, tokenService TokenService, logger *log.Logger, transactor Transactor) *Handler {
	return &Handler{
		userService:  userService,
		tokenService: tokenService,
		logger:       logger,
		transactor:   transactor,
	}
}
