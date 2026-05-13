package api

import (
	"fmt"
	"net/http"
	"todoshnik/internal/api/response"
	"todoshnik/internal/task"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	status := params.Get("status")
	scope := getScope(r)

	tasks, err := h.service.List(
		r.Context(),
		task.TaskFilter{
			Status: task.Status(status),
			Scope:  scope,
		})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("UserID: %d | Запрошены задачи\n", scope.UserID)
	response.WriteJSON(w, http.StatusOK, tasks)
}
