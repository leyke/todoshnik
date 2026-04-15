package api

import (
	"fmt"
	"log"
	"net/http"
	"todoshnik/internal/api/task"
	"todoshnik/internal/app"
)

type APIHandler struct {
	taskHandler *task.Handler
	logger      *log.Logger
}

func NewAPIHandler(container *app.App, logger *log.Logger) *APIHandler {
	return &APIHandler{
		taskHandler: task.NewHandler(container.TaskService),
		logger:      logger,
	}
}

func (api *APIHandler) Run() {
	fmt.Printf("Hello\n")
	r := api.Router()

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		fmt.Println("Error: setting up server")
	}
}

func (api *APIHandler) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
