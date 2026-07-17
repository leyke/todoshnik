package auth

import (
	"errors"
	"net/http"

	"todoshnik/cmd/api/app/request"
	"todoshnik/cmd/api/app/response"
	"todoshnik/internal/domains/token"

	apierrors "todoshnik/cmd/api/app/errors"
	tokenerror "todoshnik/internal/domains/token/errors"
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
		response.WriteError(w, apierrors.ErrNotFound)
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), user, token.DeviceTypeApi)

	if errors.Is(err, tokenerror.ErrUserNotFound) {
		response.WriteError(w, apierrors.ErrNotFound)
		return
	}

	if err != nil {
		response.WriteError(w, apierrors.ErrUnauth)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
