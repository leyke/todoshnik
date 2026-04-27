package user

import (
	"sort"
	"sync"
	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/storage"
)

type UserFileRepository struct {
	mu      sync.RWMutex
	storage storage.FileStorage[domain.User]
	users   map[int]*domain.User
	nextID  int
}

func NewUserFileRepository(storage storage.FileStorage[domain.User]) (*UserFileRepository, error) {
	tasks, err := storage.Load()
	if err != nil {
		return nil, err
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return &UserFileRepository{
		storage: storage,
		users:   tasks,
		nextID:  maxID + 1,
	}, nil
}

func (repo *UserFileRepository) List() []*domain.User {
	keys := make([]int, 0, len(repo.users))
	for k := range repo.users {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	result := make([]*domain.User, 0, len(keys))
	for _, k := range keys {
		result = append(result, repo.users[k])
	}
	return result
}

func (repo *UserFileRepository) GetByID(id int) (*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, ok := repo.users[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}

func (repo *UserFileRepository) GetUserByTgId(userTgId int64) (*domain.User, error) {
	for _, user := range repo.users {
		if user.TelegramID == userTgId {
			return user, nil
		}
	}

	return nil, apperrors.ErrNotFound
}

func (repo *UserFileRepository) Create(user *domain.User) (*domain.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user.ID = repo.nextID
	repo.nextID++

	repo.users[user.ID] = user

	err := repo.storage.Save(repo.users)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserFileRepository) Update(user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	prev := repo.users[user.ID]
	repo.users[user.ID] = user

	err := repo.storage.Save(repo.users)
	if err != nil {
		repo.users[user.ID] = prev
		return err
	}

	return nil
}

func (repo *UserFileRepository) Delete(user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	prev := user
	delete(repo.users, user.ID)

	err := repo.storage.Save(repo.users)
	if err != nil {
		repo.users[user.ID] = prev
		return err
	}

	return nil
}
