package service

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/storage"

	"todoshnik/internal/validation"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	users   map[int]*domain.User
	nextID  int
	storage *storage.FileStorage[domain.User]
	mu      sync.Mutex
}

func NewUserService() (*UserService, error) {
	storagePath := os.Getenv("tmp_dir") + "/users.json"
	storage := storage.NewFileStorage[domain.User](storagePath)
	users, err := storage.Load()
	if err != nil {
		return nil, err
	}

	s := &UserService{
		users:   users,
		storage: storage,
	}

	maxID := 0
	for _, users := range users {
		if users.ID > maxID {
			maxID = users.ID
		}
	}
	s.nextID = maxID + 1

	return s, nil
}

func (s *UserService) AddUser(name string, telegramID int64) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newUser := &domain.User{
		ID:         s.nextID,
		Name:       name,
		TelegramID: telegramID,
	}

	ve := validation.Validate(newUser)
	if ve != nil {
		fmt.Println(ve)
		return nil, apperrors.NewValidationErrorFromValidator(ve.(validator.ValidationErrors))
	}

	s.users[s.nextID] = newUser
	s.nextID++

	err := s.storage.Save(s.users)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *UserService) ListUsers() []*domain.User {
	keys := make([]int, 0, len(s.users))
	for k := range s.users {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	result := make([]*domain.User, 0, len(keys))
	for _, k := range keys {
		result = append(result, s.users[k])
	}
	return result
}

func (s *UserService) UpdateUser(userId int, name string) (*domain.User, error) {
	user, ok := s.users[userId]
	if !ok {
		return nil, apperrors.ErrNotFound
	}

	prev := user
	user.Name = name
	validateError := validation.Validate(user)
	if validateError != nil {
		user.Name = prev.Name
		return nil, apperrors.NewValidationErrorFromValidator(validateError.(validator.ValidationErrors))
	}

	err := s.storage.Save(s.users)
	if err != nil {
		s.users[userId] = prev
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(userId int) error {
	_, ok := s.users[userId]
	if !ok {
		return apperrors.ErrNotFound
	}

	delete(s.users, userId)

	return s.storage.Save(s.users)
}

func (s *UserService) GetUser(userId int) (*domain.User, error) {
	user, ok := s.users[userId]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

func (s *UserService) GetUserByTgId(userTgId int64) (*domain.User, error) {
	for _, user := range s.users {
		if user.TelegramID == userTgId {
			return user, nil
		}
	}

	return nil, apperrors.ErrNotFound
}
