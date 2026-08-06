package api

import (
	"net/http"

	"todoshnik/cmd/api/app/response"

	authcontext "todoshnik/internal/infrastructure/context_manager/auth"
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

	// на случай проверки прав или какой либо логики по юзеру, пока проверяю только наличие в базе
	_, err = h.userGetter.GetById(r.Context(), userID)
	if err != nil {
		http.Error(w, "Unauthorized: не найден пользователь", http.StatusUnauthorized)
		return
	}

	scope := getScope(userID)
	task, updateErr := h.service.MarkDone(r.Context(), id, scope)
	if updateErr != nil {
		response.WriteError(w, updateErr)
		return
	}

	h.logger.Printf("UserID: %d | Отмечена как выполненная задача: %v\n", scope.UserID, id)
	response.WriteJSON(w, http.StatusOK, task)
}
