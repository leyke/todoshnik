package api

import (
	"todoshnik/internal/task"
)

type Handler struct {
	service *task.Service
}

func NewHandler(s *task.Service) *Handler {
	return &Handler{service: s}
}
