package api

import (
	"fmt"
	"net/http"

	"todoshnik/internal/api/request"
	"todoshnik/internal/api/response"
	"todoshnik/internal/auth"

	apierrors "todoshnik/internal/api/errors"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var requestDto *UserSignInRequestDto

	requestDto, err := request.DecodeJSON[UserSignInRequestDto](r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	user, err := h.userService.GetByLogin(r.Context(), requestDto.Login)
	if err != nil {
		response.WriteError(w, fmt.Errorf("%w: неверный логин или пароль", apierrors.ErrUnauth))
		return
	}

	if !auth.ComparePassword(user.PasswordHash, requestDto.Password) {
		response.WriteError(w, fmt.Errorf("%w: неверный логин или пароль", apierrors.ErrUnauth))
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), user, auth.DeviceTypeApi)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
