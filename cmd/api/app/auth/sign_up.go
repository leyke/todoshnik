package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"todoshnik/cmd/api/app/request"
	"todoshnik/cmd/api/app/response"
	"todoshnik/internal/domains/token"
	"todoshnik/internal/infrastructure/validation"

	apierrors "todoshnik/cmd/api/app/errors"
	appuser "todoshnik/internal/domains/user"
	usererrors "todoshnik/internal/domains/user/errors"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var (
		requestDto  *UserSignUpRequestDto
		newUser     *appuser.User
		accessToken string
	)

	requestDto, err := request.DecodeJSON[UserSignUpRequestDto](r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	err = h.transactor.WithinTransaction(r.Context(), func(ctx context.Context) error {
		var err error

		newUser, err = h.userService.Add(ctx, requestDto.Name, requestDto.Login, requestDto.Password)
		if err != nil {
			return err
		}

		accessToken, err = h.tokenService.Add(ctx, newUser, token.DeviceTypeApi)
		if err != nil {
			return err
		}

		return nil
	})

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

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      newUser.ID,
		AccessToken: accessToken,
	})
}
