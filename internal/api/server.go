package api

import (
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
}

func NewAPIHandler(container *app.App) *APIHandler {
	return &APIHandler{
		taskHandler: taskapi.NewHandler(container.TaskService),
		authHandler: authapi.NewHandler(container.UserService, container.TokenService),
		logger:      container.Logger,
	}
}

func (api *APIHandler) Run() {
	fmt.Printf("Hello\n")
	r := api.Router()

	err := http.ListenAndServe(":"+os.Getenv("API_PORT"), r)
	if err != nil {
		fmt.Println("Error: setting up server")
	}
}

func (api *APIHandler) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
