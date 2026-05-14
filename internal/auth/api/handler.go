package api

import (
	"todoshnik/internal/auth/token"
	"todoshnik/internal/user"
)

type Handler struct {
	userService  *user.Service
	tokenService *token.Service
}

func NewHandler(userService *user.Service, tokenService *token.Service) *Handler {
	return &Handler{
		userService:  userService,
		tokenService: tokenService,
	}
}
