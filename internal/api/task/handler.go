package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todoshnik/internal/api/contextkeys"
	"todoshnik/internal/api/dto"
	"todoshnik/internal/api/response"
	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *service.TaskService
}

func NewHandler(s *service.TaskService) *Handler {
	return &Handler{service: s}
}

func (api *Handler) List(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	method := params.Get("status")
	userID := r.Context().Value(contextkeys.UserIDKey).(int)

	tasks := api.service.ListTasks(domain.TaskFilter{
		Status: domain.TaskStatus(method),
		Scope:  domain.AccessScope{IsAdmin: false, UserID: userID},
	})

	fmt.Printf("UserID: %d | Запрошены задачи\n", userID)
	response.WriteJSON(w, http.StatusOK, tasks)
}

func (api *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var requestDto dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(contextkeys.UserIDKey).(int)

	task, err := api.service.AddTask(requestDto.Title, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("UserID: %d | Создана задача: %v\n", userID, task.ID)
	response.WriteJSON(w, http.StatusOK, task)

}

func (api *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var requestDto dto.UpdateTaskRequest
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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

	userID := r.Context().Value(contextkeys.UserIDKey).(int)
	task, err := api.service.UpdateTask(
		id,
		requestDto.Title,
		requestDto.Done,
		domain.AccessScope{IsAdmin: false, UserID: userID},
	)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	fmt.Printf("UserID: %d | Обновлена задача: %v\n", userID, task.ID)
	response.WriteJSON(w, http.StatusOK, task)
}

func (api *Handler) Done(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(contextkeys.UserIDKey).(int)
	updateErr := api.service.MarkDone(id, domain.AccessScope{IsAdmin: false, UserID: userID})
	if updateErr != nil {
		if errors.Is(updateErr, apperrors.ErrNotFound) {
			http.Error(w, updateErr.Error(), http.StatusNotFound)
		} else {
			http.Error(w, updateErr.Error(), http.StatusBadRequest)
		}
		return
	}

	fmt.Printf("UserID: %d | Отмечена как выполненная задача: %v\n", userID, id)
	response.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (api *Handler) View(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(contextkeys.UserIDKey).(int)
	task, err := api.service.GetTask(id, domain.AccessScope{IsAdmin: false, UserID: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Printf("UserID: %d | Просмотрена задача: %v\n", userID, task.ID)
	response.WriteJSON(w, http.StatusCreated, task)
}

func (api *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(contextkeys.UserIDKey).(int)
	err = api.service.DeleteTask(id, domain.AccessScope{IsAdmin: false, UserID: userID})
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	fmt.Printf("UserID: %d | Удалена задача: %v\n", userID, id)
	response.WriteJSON(w, http.StatusNoContent, nil)
}
