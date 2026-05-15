package response

import (
	"encoding/json"
	"errors"
	"net/http"
	apperrors "todoshnik/internal/errors"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	default:
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
}
