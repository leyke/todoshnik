package user

import (
	"context"
	"fmt"

	"todoshnik/internal/auth"
	apperrors "todoshnik/internal/errors"

	"todoshnik/internal/validation"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(ctx context.Context, name string, login string, password string) (*User, error) {
	user, _ := s.repo.GetByLogin(ctx, login)
	if user != nil {
		return nil, apperrors.ErrConflict
	}

	passwordHash := auth.HashPassword(password)

	newUser := &User{
		Name:         name,
		Login:        login,
		PasswordHash: passwordHash,
	}

	ve := validateUser(newUser)
	if ve != nil {
		return nil, ve
	}

	user, err := s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) AddFromTg(ctx context.Context, name string, telegramID int64) (*User, error) {
	user, _ := s.repo.GetByTgId(ctx, telegramID)
	if user != nil {
		return user, nil
	}

	newUser := &User{
		Name:       name,
		TelegramID: telegramID,
	}

	ve := validateUser(newUser)
	if ve != nil {
		return nil, ve
	}

	user, err := s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) List(ctx context.Context) []*User {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, userID int, name string) (*User, error) {
	user, errNotFound := s.Get(ctx, userID, "")
	if errNotFound != nil {
		return nil, apperrors.ErrNotFound
	}

	prev := user
	user.Name = name
	ve := validateUser(user)
	if ve != nil {
		user.Name = prev.Name
		return nil, ve
	}

	err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Delete(ctx context.Context, userID int) error {
	user, err := s.Get(ctx, userID, "")
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, user)
}

func (s *Service) Get(ctx context.Context, userID int, login string) (*User, error) {
	var user *User
	var ok bool

	if login != "" {
		user, ok = s.repo.GetByLogin(ctx, login)
	} else if userID != 0 {
		user, ok = s.repo.GetByID(ctx, userID)
	}

	if !ok {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}

func (s *Service) GetByTgId(ctx context.Context, userTgID int64) (*User, error) {
	user, ok := s.repo.GetByTgId(ctx, userTgID)
	if !ok {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}

func validateUser(user *User) error {
	ve := validation.Validate(user)
	if ve != nil {
		fmt.Println(ve)
		return apperrors.NewValidationErrorFromValidator(ve)
	}
	return nil
}
