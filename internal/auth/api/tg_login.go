package api

import (
	"net/http"
	"todoshnik/internal/api/request"
	"todoshnik/internal/api/response"
	"todoshnik/internal/auth"
)

func (h *Handler) TgLogin(w http.ResponseWriter, r *http.Request) {
	var requestDto *TgLoginRequestDto
	requestDto, err := request.DecodeJSON[TgLoginRequestDto](r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	user, err := h.userService.GetByTgId(r.Context(), requestDto.TgUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), user, auth.DeviceTypeBot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// TODO может ли пользователь удалиться после получения userService.GetByTgId и перед tokenService.Add и что будет?

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
