package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"todoshnik/internal/api/task"
	"todoshnik/internal/app"
)

type APIHandler struct {
	taskHandler *task.Handler
	logger      *log.Logger
}

func NewAPIHandler(container *app.App) *APIHandler {
	return &APIHandler{
		taskHandler: task.NewHandler(container.TaskService),
		logger:      container.Logger,
	}
}

func (api *APIHandler) Run() {
	fmt.Printf("Hello\n")
	r := api.Router()

	err := http.ListenAndServe(os.Getenv("host")+":"+os.Getenv("port"), r)
	if err != nil {
		fmt.Println("Error: setting up server")
	}
}

func (api *APIHandler) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
