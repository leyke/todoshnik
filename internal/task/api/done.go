package api

import (
	"fmt"
	"net/http"
	"todoshnik/internal/api/response"
)

func (h *Handler) Done(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskID(r)
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	scope := getScope(r)
	task, updateErr := h.service.MarkDone(r.Context(), id, scope)
	if updateErr != nil {
		writeError(w, updateErr)
		return
	}

	fmt.Printf("UserID: %d | Отмечена как выполненная задача: %v\n", scope.UserID, id)
	response.WriteJSON(w, http.StatusOK, task)
}
