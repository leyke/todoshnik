package user

import (
	"context"
	"fmt"

	"todoshnik/internal/infrastructure/validation"

	usererrors "todoshnik/internal/domains/user/errors"
)

type Service struct {
	repo           Repository
	passwordHasher PasswordHasher
}

func NewService(repo Repository, passwordHasher PasswordHasher) *Service {
	return &Service{
		repo:           repo,
		passwordHasher: passwordHasher,
	}
}

func (s *Service) Add(ctx context.Context, name string, login string, password string) (*User, error) {
	user, _ := s.repo.GetByLogin(ctx, login)
	if user != nil {
		return nil, usererrors.ErrConflict
	}

	passwordHash, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		Name:         name,
		Login:        login,
		PasswordHash: passwordHash,
	}

	ve := validateUser(newUser)
	if ve != nil {
		return nil, ve
	}

	user, err = s.repo.Create(ctx, newUser)
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

func (s *Service) List(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, userID int, name string) (*User, error) {
	user, errNotFound := s.GetById(ctx, userID)
	if errNotFound != nil {
		return nil, usererrors.ErrNotFound
	}

	updated := *user
	updated.Name = name

	if err := validateUser(&updated); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (s *Service) Delete(ctx context.Context, userID int) error {
	user, err := s.GetById(ctx, userID)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, user)
}

func (s *Service) GetByLogin(ctx context.Context, login string) (*User, error) {
	user, err := s.repo.GetByLogin(ctx, login)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) ValidatePassword(hash string, password string) (bool, error) {
	return s.passwordHasher.Compare(hash, password)
}

func (s *Service) GetById(ctx context.Context, userID int) (*User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByTgId(ctx context.Context, userTgID int64) (*User, error) {
	user, err := s.repo.GetByTgId(ctx, userTgID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func validateUser(user *User) error {
	ve := validation.Validate(user)
	if ve != nil {
		fmt.Println(ve)
		return validation.NewValidationErrorFromValidator(ve)
	}
	return nil
}
