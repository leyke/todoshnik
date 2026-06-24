package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"todoshnik/internal/api/response"
	"todoshnik/internal/app"
	"todoshnik/internal/config"

	authapi "todoshnik/internal/auth/api"
	taskapi "todoshnik/internal/domains/task/api"
)

type APIHandler struct {
	taskHandler *taskapi.Handler
	authHandler *authapi.Handler
	logger      *log.Logger
	server      *http.Server

	config Config
}

func NewAPIHandler(services *app.Services, logger *log.Logger, config Config) *APIHandler {
	return &APIHandler{
		taskHandler: taskapi.NewHandler(services.TaskService),
		authHandler: authapi.NewHandler(services.UserService, services.TokenService, logger),
		logger:      logger,
		config:      config,
	}
}

func (h *APIHandler) Run(cfg config.AppConfig) error {
	h.server = &http.Server{
		Addr:    ":" + cfg.Port, // в докере с хостом не работает
		Handler: h.Router(),
	}

	h.logger.Println("API hello")

	err := h.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("API server error: %w", err)
	}

	return nil
}

func (h *APIHandler) Shutdown(ctx context.Context) error {
	h.logger.Println("отключение API...")
	return h.server.Shutdown(ctx)
}

func (h *APIHandler) pingHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, "pong")
}
