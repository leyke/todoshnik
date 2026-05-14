package api

import (
	"encoding/json"
	"net/http"
	"todoshnik/internal/api/response"
	"todoshnik/internal/auth"
)

func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var requestDto UserSignInRequestDto
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Get(r.Context(), 0, requestDto.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !auth.ComparePassword(user.PasswordHash, requestDto.Password) {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), user, auth.DeviceTypeApi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
