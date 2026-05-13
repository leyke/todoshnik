package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todoshnik/internal/api/response"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var requestDto UpdateTaskRequest
	id, err := getTaskID(r)
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	scope := getScope(r)
	task, err := h.service.Update(
		r.Context(),
		id,
		requestDto.Title,
		requestDto.Done,
		scope,
	)

	if err != nil {
		writeError(w, err)
		return
	}

	fmt.Printf("UserID: %d | Обновлена задача: %v\n", scope.UserID, task.ID)
	response.WriteJSON(w, http.StatusOK, task)
}
