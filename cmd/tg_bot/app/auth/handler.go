package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	client "todoshnik/internal/infrastructure/api_client"
)

type Api interface {
	Post(ctx context.Context, endpoint string, payload any) (*http.Response, error)
}

type Config struct {
	RequestTimeout time.Duration
}

type Handler struct {
	api            *client.ApiClient
	requestTimeout time.Duration
	logger         *log.Logger
}

func NewHandler(api *client.ApiClient, config Config, l *log.Logger) *Handler {
	return &Handler{
		requestTimeout: config.RequestTimeout,
		api:            api,
		logger:         l,
	}
}
