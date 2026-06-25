package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	apierrors "todoshnik/cmd/api/app/errors"
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
	case errors.Is(err, apierrors.ErrNotFound):
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apierrors.ErrEmptyBody):
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apierrors.ErrBadRequest):
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apierrors.ErrUnauth):
		WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	case errors.Is(err, apierrors.ErrConflict):
		WriteJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
	default:
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}
}
