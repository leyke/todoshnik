package api

import (
	"fmt"
	"net/http"
	"todoshnik/internal/api/response"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskID(r)
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	scope := getScope(r)
	err = h.service.Delete(r.Context(), id, scope)
	if err != nil {
		writeError(w, err)
		return
	}

	fmt.Printf("UserID: %d | Удалена задача: %v\n", scope.UserID, id)
	response.WriteJSON(w, http.StatusNoContent, nil)
}
