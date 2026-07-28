package api

import (
	"fmt"
	"net/http"

	"todoshnik/cmd/api/app/response"

	authcontext "todoshnik/internal/infrastructure/context_manager/auth"
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

	// на случай проверки прав или какой либо логики по юзеру, пока проверяю только наличие в базе
	_, err = h.userGetter.GetById(r.Context(), userID)
	if err != nil {
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
