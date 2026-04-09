package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todoshnik/internal/helpers"
	"todoshnik/internal/service"
)

type APIHandler struct {
	service *service.TaskService
}

func NewAPIHandler(s *service.TaskService) *APIHandler {
	return &APIHandler{service: s}
}

func (api *APIHandler) Run() {
	fmt.Printf("Hello\n")

	http.HandleFunc("/ping", api.pingHandler)
	http.HandleFunc("/tasks", api.tasksHandler)

	http.ListenAndServe(":8080", nil)
}

func (api *APIHandler) tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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

		return
	}
	params := r.URL.Query()
	method := params.Get("status")
	tasks := api.service.ListTasks(method)
	if len(tasks) == 0 {
		fmt.Fprintf(w, "Список задач пуст")
	}
	fmt.Printf("Запрошены задачи\n")
	json.NewEncoder(w).Encode(tasks)
}

func (api *APIHandler) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
