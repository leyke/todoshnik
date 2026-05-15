package api

import (
	"fmt"
	"net/http"
	"todoshnik/internal/api/response"
	authcontext "todoshnik/internal/auth/context"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
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
	err = h.service.Delete(r.Context(), id, scope)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	fmt.Printf("UserID: %d | Удалена задача: %v\n", scope.UserID, id)
	response.WriteJSON(w, http.StatusNoContent, nil)
}
