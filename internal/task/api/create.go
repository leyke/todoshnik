package api

import (
	"fmt"
	"net/http"
	"todoshnik/internal/api/request"
	"todoshnik/internal/api/response"
	authcontext "todoshnik/internal/auth/context"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var requestDto *CreateTaskRequest
	requestDto, err := request.DecodeJSON[CreateTaskRequest](r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	userID, ok := authcontext.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: не найден пользователь", http.StatusUnauthorized)
		return
	}

	task, err := h.service.Add(r.Context(), requestDto.Title, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	fmt.Printf("UserID: %d | Создана задача: %v\n", userID, task.ID)
	response.WriteJSON(w, http.StatusOK, task)

}
