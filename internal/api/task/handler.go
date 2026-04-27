package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

	// TODO: userid из AuthMiddleware
	tasks := api.service.ListTasks(domain.TaskFilter{
		Status: domain.TaskStatus(method),
		Scope:  domain.AccessScope{IsAdmin: true},
	})

	fmt.Printf("Запрошены задачи\n")
	response.WriteJSON(w, http.StatusOK, tasks)
}

func (api *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var requestDto dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// TODO: userid из AuthMiddleware
	userID, err := strconv.Atoi(requestDto.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := api.service.AddTask(requestDto.Title, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Создана задача: %v\n", task.ID)
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
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	// TODO: userid из AuthMiddleware
	task, err := api.service.UpdateTask(id, requestDto.Title, requestDto.Done, domain.AccessScope{IsAdmin: true})
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	fmt.Printf("Обновлена задача: %v\n", task.ID)
	response.WriteJSON(w, http.StatusOK, task)
}

func (api *Handler) Done(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	// TODO: userid из AuthMiddleware
	updateErr := api.service.MarkDone(id, domain.AccessScope{IsAdmin: true})
	if updateErr != nil {
		if errors.Is(updateErr, apperrors.ErrNotFound) {
			http.Error(w, updateErr.Error(), http.StatusNotFound)
		} else {
			http.Error(w, updateErr.Error(), http.StatusBadRequest)
		}
		return
	}

	fmt.Printf("Отмечена как выполненная задача: %v\n", id)
	response.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (api *Handler) View(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	// TODO: userid из AuthMiddleware
	task, err := api.service.GetTask(id, domain.AccessScope{IsAdmin: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Printf("Просмотрена задача: %v\n", task.ID)
	response.WriteJSON(w, http.StatusCreated, task)
}

func (api *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	// TODO: userid из AuthMiddleware
	err = api.service.DeleteTask(id, domain.AccessScope{IsAdmin: true})
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	fmt.Printf("Удалена задача: %v\n", id)
	response.WriteJSON(w, http.StatusNoContent, nil)
}
