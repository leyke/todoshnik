package auth

import (
	"log"
	
	"todoshnik/internal/domains/token"
	"todoshnik/internal/domains/user"
	"todoshnik/internal/infrastructure/db/transaction"
)

type Handler struct {
	userService  *user.Service
	tokenService *token.Service
	logger       *log.Logger
	transactor   *transaction.Transactor
}

func NewHandler(userService *user.Service, tokenService *token.Service, logger *log.Logger, transactor *transaction.Transactor) *Handler {
	return &Handler{
		userService:  userService,
		tokenService: tokenService,
		logger:       logger,
		transactor:   transactor,
	}
}
