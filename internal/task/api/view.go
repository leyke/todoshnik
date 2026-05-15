package api

import (
	"fmt"
	"net/http"
	"todoshnik/internal/api/response"
	authcontext "todoshnik/internal/auth/context"
)

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
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
	task, err := h.service.Get(r.Context(), id, scope)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Printf("UserID: %d | Просмотрена задача: %v\n", scope.UserID, task.ID)
	response.WriteJSON(w, http.StatusOK, task)
}
