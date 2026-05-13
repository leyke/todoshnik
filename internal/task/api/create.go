package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todoshnik/internal/api/response"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var requestDto CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	userID := getUserID(r)

	task, err := h.service.Add(r.Context(), requestDto.Title, userID)
	if err != nil {
		writeError(w, err)
		return
	}
	fmt.Printf("UserID: %d | Создана задача: %v\n", userID, task.ID)
	response.WriteJSON(w, http.StatusOK, task)

}
