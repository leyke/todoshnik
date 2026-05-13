package auth

import (
	"context"
	"encoding/json"
	"net/http"
	authapi "todoshnik/internal/auth/api"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/user"

	"todoshnik/internal/api/response"
	"todoshnik/internal/auth"
	"todoshnik/internal/domain"
	"todoshnik/internal/service"
)

type AuthHandler struct {
	userService  *user.Service
	tokenService *service.AccessTokenService
}

func NewAuthHandler(userService *user.Service, tokenService *service.AccessTokenService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (h AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var requestDto authapi.UserSignUpRequestDto
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
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

	accessToken, err := h.tokenService.AddToken(r.Context(), user, domain.DeviceTypeApi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.WriteJSON(w, http.StatusOK, authapi.AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken.Hash,
	})
}

func (h AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var requestDto authapi.UserSignInRequestDto
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Get(r.Context(), 0, requestDto.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !auth.ComparePassword(user.PasswordHash, requestDto.Password) {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.tokenService.AddToken(r.Context(), user, domain.DeviceTypeApi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.WriteJSON(w, http.StatusOK, authapi.AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken.Hash,
	})
}

func (h AuthHandler) TgLogin(w http.ResponseWriter, r *http.Request) {
	var requestDto authapi.TgLoginRequestDto
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByTgId(r.Context(), requestDto.TgUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	accessToken, err := h.tokenService.AddToken(r.Context(), user, domain.DeviceTypeBot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.WriteJSON(w, http.StatusOK, authapi.AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken.Hash,
	})
}

func (h AuthHandler) TgAutoReg(w http.ResponseWriter, r *http.Request) {
	var requestDto authapi.TgLoginRequestDto
	if err := json.NewDecoder(r.Body).Decode(&requestDto); err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Пустой запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.userService.AddFromTg(r.Context(), requestDto.Name, requestDto.TgUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, err := h.tokenService.AddToken(r.Context(), user, domain.DeviceTypeBot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response.WriteJSON(w, http.StatusOK, authapi.AuthResponseDto{
		UserID:      user.ID,
		AccessToken: accessToken.Hash,
	})
}

func (h AuthHandler) ValidateToken(ctx context.Context, token string) (*user.User, error) {
	userID, err := h.tokenService.GetUserID(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.Get(ctx, userID, "")
	if err != nil {
		return nil, err
	}

	return user, nil
}
