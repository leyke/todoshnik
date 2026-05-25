package bot

import (
	"log"
	"todoshnik/internal/client"
	taskgrpc "todoshnik/internal/task/grpc"
)

type Handler struct {
	api        *client.ApiClient
	taskClient *taskgrpc.Client
	logger     *log.Logger
}

func NewHandler(api *client.ApiClient, taskClient *taskgrpc.Client, l *log.Logger) *Handler {
	return &Handler{api: api, taskClient: taskClient, logger: l}
}
