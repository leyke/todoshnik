package api

import (
	"context"
	"todoshnik/internal/domains/task"
	"todoshnik/internal/domains/user"
)

type UserGetter interface {
	GetById(ctx context.Context, userID int) (*user.User, error)
}

type Handler struct {
	service    *task.Service
	userGetter UserGetter
}

func NewHandler(s *task.Service, ug UserGetter) *Handler {
	return &Handler{service: s, userGetter: ug}
}
