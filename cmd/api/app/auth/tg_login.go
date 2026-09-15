package auth

import (
	"errors"
	"net/http"

	"todoshnik/cmd/api/app/request"
	"todoshnik/cmd/api/app/response"
	"todoshnik/internal/domains/token"

	apierrors "todoshnik/cmd/api/app/errors"
	usererrors "todoshnik/internal/domains/user/errors"
)

func (h *Handler) TgLogin(w http.ResponseWriter, r *http.Request) {
	var requestDto *TgLoginRequestDto
	requestDto, err := request.DecodeJSON[TgLoginRequestDto](r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	user, err := h.userService.GetByTgId(r.Context(), requestDto.TgUserID)

	if errors.Is(err, usererrors.ErrNotFound) {
		response.WriteError(w, apierrors.ErrNotFound)
		return
	}

	if err != nil {
		h.logger.Printf("ошибка поиска пользователя: %v", err)
		response.WriteError(w, err)
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), user, token.DeviceTypeBot)

	if err != nil {
		h.logger.Printf("ошибка создания токена: %v", err)
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
