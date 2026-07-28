package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"todoshnik/cmd/api/app/response"
	"todoshnik/internal/app"
	"todoshnik/internal/config"
	"todoshnik/internal/infrastructure/db/transaction"

	authapi "todoshnik/cmd/api/app/auth"
	taskapi "todoshnik/cmd/api/app/task"
)

type APIHandler struct {
	taskHandler *taskapi.Handler
	authHandler *authapi.Handler
	logger      *log.Logger
	server      *http.Server

	config Config
}

func NewAPIHandler(services *app.Services, logger *log.Logger, transactor *transaction.Transactor, config Config) *APIHandler {
	return &APIHandler{
		taskHandler: taskapi.NewHandler(services.TaskService, services.UserService),
		authHandler: authapi.NewHandler(services.UserService, services.TokenService, logger, transactor),
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
