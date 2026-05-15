package bot

import (
	"todoshnik/internal/client"
)

type Handler struct {
	api *client.ApiClient
}

func NewHandler(api *client.ApiClient) *Handler {
	return &Handler{api: api}
}
