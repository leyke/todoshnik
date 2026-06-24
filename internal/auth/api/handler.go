package api

import (
	"log"
	"todoshnik/internal/domains/token"
	"todoshnik/internal/domains/user"
)

type Handler struct {
	userService  *user.Service
	tokenService *token.Service
	logger       *log.Logger
}

func NewHandler(userService *user.Service, tokenService *token.Service, logger *log.Logger) *Handler {
	return &Handler{
		userService:  userService,
		tokenService: tokenService,
		logger:       logger,
	}
}
