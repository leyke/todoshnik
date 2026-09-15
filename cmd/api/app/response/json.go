package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"todoshnik/internal/infrastructure/validation"

	apierrors "todoshnik/cmd/api/app/errors"
	taskerrors "todoshnik/internal/domains/task/errors"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при кодировании JSON: %v", err), http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, taskerrors.ErrNotFound),
		errors.Is(err, apierrors.ErrNotFound):
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apierrors.ErrEmptyBody):
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apierrors.ErrBadRequest):
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apierrors.ErrUnauth):
		WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apierrors.ErrConflict):
		WriteJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
	case errors.Is(err, validation.ErrNotValidate):
		WriteJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Error: err.Error()})
	default:
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}
}
