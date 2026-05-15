package api

import (
	"fmt"
	"net/http"
	"todoshnik/internal/api/response"
	authcontext "todoshnik/internal/auth/context"
	"todoshnik/internal/task"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	status := params.Get("status")

	userID, ok := authcontext.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: не найден пользователь", http.StatusUnauthorized)
		return
	}

	scope := getScope(userID)

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
