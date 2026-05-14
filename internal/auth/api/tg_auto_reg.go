package api

import (
	"encoding/json"
	"net/http"
	"todoshnik/internal/api/response"
	"todoshnik/internal/auth"
)

func (h Handler) TgAutoReg(w http.ResponseWriter, r *http.Request) {
	var requestDto TgLoginRequestDto
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.userService.AddFromTg(r.Context(), requestDto.Name, requestDto.TgUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), user, auth.DeviceTypeBot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
