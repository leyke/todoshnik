package service

import (
	"fmt"

	"todoshnik/internal/auth"
	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	repository "todoshnik/internal/repository/user"

	"todoshnik/internal/validation"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) AddUser(name string, login string, password string) (*domain.User, error) {
	user, _ := s.repo.GetByLogin(login)
	if user != nil {
		return nil, apperrors.ErrConflict
	}

	passwordHash := auth.HashPassword(password)

	newUser := &domain.User{
		Name:         name,
		Login:        login,
		PasswordHash: passwordHash,
	}

	ve := validation.Validate(newUser)
	if ve != nil {
		fmt.Println(ve)
		return nil, apperrors.NewValidationErrorFromValidator(ve.(validator.ValidationErrors))
	}

	user, err := s.repo.Create(newUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) AddTgUser(name string, telegramID int64) (*domain.User, error) {
	user, _ := s.repo.GetUserByTgId(telegramID)
	if user != nil {
		return user, nil
	}

	newUser := &domain.User{
		Name:       name,
		TelegramID: telegramID,
	}

	ve := validation.Validate(newUser)
	if ve != nil {
		fmt.Println(ve)
		return nil, apperrors.NewValidationErrorFromValidator(ve.(validator.ValidationErrors))
	}

	user, err := s.repo.Create(newUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ListUsers() []*domain.User {
	return s.repo.List()
}

func (s *UserService) UpdateUser(userID int, name string) (*domain.User, error) {
	user, errNotFound := s.GetUser(userID, "")
	if errNotFound != nil {
		return nil, apperrors.ErrNotFound
	}

	prev := user
	user.Name = name
	validateError := validation.Validate(user)
	if validateError != nil {
		user.Name = prev.Name
		return nil, apperrors.NewValidationErrorFromValidator(validateError.(validator.ValidationErrors))
	}

	err := s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(userID int) error {
	user, err := s.GetUser(userID, "")
	if err != nil {
		return err
	}

	return s.repo.Delete(user)
}

func (s *UserService) GetUser(userID int, login string) (*domain.User, error) {
	var user *domain.User
	var ok bool

	if login != "" {
		user, ok = s.repo.GetByLogin(login)
	} else if userID != 0 {
		user, ok = s.repo.GetByID(userID)
	}

	if !ok {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}

func (s *UserService) GetUserByTgId(userTgID int64) (*domain.User, error) {
	user, ok := s.repo.GetUserByTgId(userTgID)
	if !ok {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}
