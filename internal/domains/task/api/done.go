package api

import (
	"fmt"
	"net/http"

	"todoshnik/internal/api/response"

	authcontext "todoshnik/internal/auth/context"
)

func (h *Handler) Done(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskID(r)
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	userID, ok := authcontext.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: не найден пользователь", http.StatusUnauthorized)
		return
	}

	scope := getScope(userID)
	task, updateErr := h.service.MarkDone(r.Context(), id, scope)
	if updateErr != nil {
		response.WriteError(w, updateErr)
		return
	}

	fmt.Printf("UserID: %d | Отмечена как выполненная задача: %v\n", scope.UserID, id)
	response.WriteJSON(w, http.StatusOK, task)
}
