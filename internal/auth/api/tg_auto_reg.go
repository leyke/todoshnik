package api

import (
	"net/http"
	"todoshnik/internal/api/request"
	"todoshnik/internal/api/response"
	"todoshnik/internal/auth"
)

func (h *Handler) TgAutoReg(w http.ResponseWriter, r *http.Request) {
	var requestDto *TgLoginRequestDto
	requestDto, err := request.DecodeJSON[TgLoginRequestDto](r)
	if err != nil {
		response.WriteError(w, err)
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
	// TODO что произойдет если мы выполним AddFromTg и сфейлимся на tokenService.Add?

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
