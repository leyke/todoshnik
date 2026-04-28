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
	items   map[int]*domain.User
	nextID  int
}

func NewUserFileRepository(storage storage.FileStorage[domain.User]) (*UserFileRepository, error) {
	items, err := storage.Load()
	if err != nil {
		return nil, err
	}

	maxID := 0
	for _, item := range items {
		if item.ID > maxID {
			maxID = item.ID
		}
	}

	return &UserFileRepository{
		storage: storage,
		items:   items,
		nextID:  maxID + 1,
	}, nil
}

func (repo *UserFileRepository) List() []*domain.User {
	keys := make([]int, 0, len(repo.items))
	for k := range repo.items {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	result := make([]*domain.User, 0, len(keys))
	for _, k := range keys {
		result = append(result, repo.items[k])
	}
	return result
}

func (repo *UserFileRepository) GetByID(id int) (*domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, ok := repo.items[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}

func (repo *UserFileRepository) GetUserByTgId(userTgId int64) (*domain.User, error) {
	for _, user := range repo.items {
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

	repo.items[user.ID] = user

	err := repo.storage.Save(repo.items)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserFileRepository) Update(user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	prev := repo.items[user.ID]
	repo.items[user.ID] = user

	err := repo.storage.Save(repo.items)
	if err != nil {
		repo.items[user.ID] = prev
		return err
	}

	return nil
}

func (repo *UserFileRepository) Delete(user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	prev := user
	delete(repo.items, user.ID)

	err := repo.storage.Save(repo.items)
	if err != nil {
		repo.items[user.ID] = prev
		return err
	}

	return nil
}
