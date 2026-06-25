package api

import (
	"fmt"
	"net/http"

	"todoshnik/cmd/api/app/request"
	"todoshnik/cmd/api/app/response"

	authcontext "todoshnik/internal/infrastructure/context_manager/auth"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var requestDto *UpdateTaskRequest
	id, err := getTaskID(r)
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	requestDto, err = request.DecodeJSON[UpdateTaskRequest](r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	userID, ok := authcontext.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: не найден пользователь", http.StatusUnauthorized)
		return
	}

	scope := getScope(userID)
	task, err := h.service.Update(
		r.Context(),
		id,
		requestDto.Title,
		requestDto.Done,
		scope,
	)

	if err != nil {
		response.WriteError(w, err)
		return
	}

	fmt.Printf("UserID: %d | Обновлена задача: %v\n", scope.UserID, task.ID)
	response.WriteJSON(w, http.StatusOK, task)
}
