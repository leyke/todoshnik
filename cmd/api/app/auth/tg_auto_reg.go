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
)

func (h *Handler) TgAutoReg(w http.ResponseWriter, r *http.Request) {
	var (
		requestDto  *TgLoginRequestDto
		user        *appuser.User
		accessToken string
	)

	requestDto, err := request.DecodeJSON[TgLoginRequestDto](r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	err = h.transactor.WithinTransaction(r.Context(), func(ctx context.Context) error {
		var err error

		user, err = h.userService.AddFromTg(ctx, requestDto.Name, requestDto.TgUserID)
		if err != nil {
			return err
		}

		accessToken, err = h.tokenService.Add(ctx, user, token.DeviceTypeBot)
		if err != nil {
			return err
		}
		return nil
	})

	if errors.Is(err, validation.ErrNotValidate) {
		response.WriteError(w, fmt.Errorf("%w: данные пользователя недействительны: %w", apierrors.ErrBadRequest, err))
		return
	}
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
