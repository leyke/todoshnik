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

	// опять нет проверки ошибок или подавления
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	// я бы сказал что StatusBadRequest это клиентская ошибка и довольно неожиданно обвинять клиента в что что-то
	// в приложении не так, в пакете error у тебя более богатый спектр ошибок
	default:
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
}
