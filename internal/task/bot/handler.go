package bot

import (
	"log"

	"todoshnik/internal/client"
)

const taskAPIURL string = "/tasks/"

type Handler struct {
	api    *client.ApiClient
	logger *log.Logger
}

func NewHandler(api *client.ApiClient, l *log.Logger) *Handler {
	return &Handler{api: api, logger: l}
}
