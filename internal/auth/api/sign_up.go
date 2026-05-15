package api

import (
	"net/http"
	"todoshnik/internal/api/request"
	"todoshnik/internal/api/response"
	"todoshnik/internal/auth"
	apperrors "todoshnik/internal/errors"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var requestDto *UserSignUpRequestDto
	requestDto, err := request.DecodeJSON[UserSignUpRequestDto](r)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	user, err := h.userService.Add(r.Context(), requestDto.Name, requestDto.Login, requestDto.Password)
	if err != nil {
		if err == apperrors.ErrConflict {
			http.Error(w, "Пользователь с таким логином уже существует", http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, err := h.tokenService.Add(r.Context(), user, auth.DeviceTypeApi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken,
	})
}
