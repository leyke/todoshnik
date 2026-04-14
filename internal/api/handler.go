package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"todoshnik/internal/helpers"
	"todoshnik/internal/service"
)

type APIHandler struct {
	service *service.TaskService
	logger  *log.Logger
}

func NewAPIHandler(s *service.TaskService, logFile *os.File) *APIHandler {
	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &APIHandler{service: s, logger: logger}
}

func (api *APIHandler) Close() {
	if api.logger != nil {
		if file, ok := api.logger.Writer().(*os.File); ok {
			file.Close()
		}
	}
}

func (api *APIHandler) Run() {
	fmt.Printf("Hello\n")

	http.Handle("/ping", api.loggingMiddleware(http.HandlerFunc(api.pingHandler)))   // --> adding middleware to route
	http.Handle("/tasks", api.loggingMiddleware(http.HandlerFunc(api.tasksHandler))) // --> adding middleware to route

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error: setting up server")
	}
}

func (api *APIHandler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.logger.Println(2, r.URL.Path)
		next.ServeHTTP(w, r)
		fmt.Println("Конец обработки")
	})
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
