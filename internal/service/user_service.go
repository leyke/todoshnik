package service

import (
	"fmt"
	"os"

	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	repo "todoshnik/internal/repository/user"
	repository "todoshnik/internal/repository/user"
	"todoshnik/internal/storage"

	"todoshnik/internal/validation"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	repo repository.UserRepositoryInreface
}

func NewUserService() (*UserService, error) {
	storagePath := os.Getenv("tmp_dir") + "/users.json"
	storage := storage.NewFileStorage[domain.User](storagePath)

	repo, err := repo.NewUserFileRepository(storage)
	if err != nil {
		return nil, err
	}

	s := &UserService{
		repo: repo,
	}

	return s, nil
}

func (s *UserService) AddUser(name string, telegramID int64) (*domain.User, error) {
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
	user, errNotFound := s.GetUser(userID)
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
	user, err := s.GetUser(userID)
	if err != nil {
		return err
	}

	return s.repo.Delete(user)
}

func (s *UserService) GetUser(userID int) (*domain.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByTgId(userTgID int64) (*domain.User, error) {
	user, err := s.repo.GetUserByTgId(userTgID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
