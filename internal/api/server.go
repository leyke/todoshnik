package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"todoshnik/internal/app"
	authapi "todoshnik/internal/auth/api"
	taskapi "todoshnik/internal/task/api"
)

type APIHandler struct {
	taskHandler *taskapi.Handler
	authHandler *authapi.Handler
	logger      *log.Logger
	server      *http.Server
}

func NewAPIHandler(container *app.App) *APIHandler {
	return &APIHandler{
		taskHandler: taskapi.NewHandler(container.TaskService),
		authHandler: authapi.NewHandler(container.UserService, container.TokenService),
		logger:      container.Logger,
	}
}

func (h *APIHandler) Run() error {
	h.server = &http.Server{
		Addr:    ":" + os.Getenv("API_PORT"),
		Handler: h.Router(),
	}

	h.logger.Println("API hello")

	err := h.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (h *APIHandler) Shutdown(ctx context.Context) error {
	h.logger.Println("отключение API...")
	return h.server.Shutdown(ctx)
}

func (h *APIHandler) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
