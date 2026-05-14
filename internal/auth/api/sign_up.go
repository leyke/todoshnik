package api

import (
	"encoding/json"
	"net/http"
	"todoshnik/internal/api/response"
	"todoshnik/internal/auth"
	apperrors "todoshnik/internal/errors"
)

func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var requestDto UserSignUpRequestDto
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Add(r.Context(), requestDto.Name, requestDto.Login, requestDto.Password)
	if err != nil {
		if err == apperrors.ErrConflict {
			http.Error(w, "Пользователь с таким логином уже существует", http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), user, auth.DeviceTypeApi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
