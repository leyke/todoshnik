package api

import (
	"todoshnik/internal/domains/task"
)

type Handler struct {
	service *task.Service
}

func NewHandler(s *task.Service) *Handler {
	return &Handler{service: s}
}
