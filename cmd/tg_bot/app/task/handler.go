package bot

import (
	"context"
	"log"
	"net/http"
	"net/url"
)

const taskAPIURL string = "/tasks/"

type Api interface {
	Get(ctx context.Context, endpoint string, query url.Values) (*http.Response, error)
	Post(ctx context.Context, endpoint string, payload any) (*http.Response, error)
	Delete(ctx context.Context, endpoint string) (*http.Response, error)
}

type Handler struct {
	api    Api
	logger *log.Logger
}

func NewHandler(api Api, l *log.Logger) *Handler {
	return &Handler{api: api, logger: l}
}
