package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todoshnik/internal/helpers"
	"todoshnik/internal/service"
)

type Handler struct {
	service *service.TaskService
}

func NewHandler(s *service.TaskService) *Handler {
	return &Handler{service: s}
}

func (api *Handler) List(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	method := params.Get("status")
	tasks := api.service.ListTasks(method)
	if len(tasks) == 0 {
		fmt.Fprintf(w, "Список задач пуст")
	}
	fmt.Printf("Запрошены задачи\n")
	json.NewEncoder(w).Encode(tasks)
}

func (api *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	title := data["title"]
	task, err := api.service.AddTask(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf(helpers.TaskPrettyPrintTemplate(), task.ID, task.Title, task.Done)
}
