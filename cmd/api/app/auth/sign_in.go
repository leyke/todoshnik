package auth

import (
	"fmt"
	"net/http"

	"todoshnik/cmd/api/app/request"
	"todoshnik/cmd/api/app/response"
	"todoshnik/internal/domains/token"

	apierrors "todoshnik/cmd/api/app/errors"
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

	isMatch, err := h.userService.ValidatePassword(user.PasswordHash, requestDto.Password)

	if err != nil {
		h.logger.Printf("ошибка проверки пароля: %v", err)
		response.WriteError(w, fmt.Errorf("%w: ошибка проверки пароля", apierrors.ErrUnauth))
		return
	}

	if !isMatch {
		response.WriteError(w, fmt.Errorf("%w: неверный логин или пароль", apierrors.ErrUnauth))
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), user, token.DeviceTypeApi)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
