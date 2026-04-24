package user

import (
	"todoshnik/internal/domain"
	"todoshnik/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	service *service.UserService
}

func NewHandler(s *service.UserService) *Handler {
	return &Handler{service: s}
}

func (uh Handler) AddUser(user *tgbotapi.User) {
	uh.service.AddUser(user.UserName, user.ID)
}

func (uh Handler) GetAppUser(tgUser *tgbotapi.User) (*domain.User, error) {
	appUser, err := uh.service.GetUserByTgId(tgUser.ID)
	return appUser, err
}
