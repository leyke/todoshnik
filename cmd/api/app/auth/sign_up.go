package auth

import (
	"errors"
	"fmt"
	"net/http"

	"todoshnik/cmd/api/app/request"
	"todoshnik/cmd/api/app/response"
	"todoshnik/internal/domains/token"
	"todoshnik/internal/infrastructure/validation"

	apierrors "todoshnik/cmd/api/app/errors"
	usererrors "todoshnik/internal/domains/user/errors"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var requestDto *UserSignUpRequestDto
	requestDto, err := request.DecodeJSON[UserSignUpRequestDto](r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	newUser, err := h.userService.Add(r.Context(), requestDto.Name, requestDto.Login, requestDto.Password)
	if errors.Is(err, usererrors.ErrConflict) {
		response.WriteError(w, fmt.Errorf("%w: пользователь с таким логином уже существует", apierrors.ErrConflict))
		return
	}
	if errors.Is(err, validation.ErrNotValidate) {
		response.WriteError(w, fmt.Errorf("%w: данные пользователя недействительны: %w", apierrors.ErrBadRequest, err))
		return
	}
	if err != nil {
		response.WriteError(w, err)
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), newUser, token.DeviceTypeApi)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      newUser.ID,
		AccessToken: accessToken,
	})
}
