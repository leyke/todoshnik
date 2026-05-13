package api

import (
	"fmt"
	"net/http"
	"todoshnik/internal/api/response"
)

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskID(r)
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	scope := getScope(r)
	task, err := h.service.Get(r.Context(), id, scope)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Printf("UserID: %d | Просмотрена задача: %v\n", scope.UserID, task.ID)
	response.WriteJSON(w, http.StatusOK, task)
}
