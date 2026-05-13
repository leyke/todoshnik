package user

import (
	"context"
	"sort"
	"sync"
	"todoshnik/internal/domain"
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

func (repo *UserFileRepository) List(ctx context.Context) []*domain.User {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

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

func (repo *UserFileRepository) GetByID(ctx context.Context, id int) (*domain.User, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, false
	default:
	}

	user, ok := repo.items[id]
	return user, ok
}

func (repo *UserFileRepository) GetByLogin(ctx context.Context, login string) (*domain.User, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, false
	default:
	}

	for _, user := range repo.items {
		if user.Login == login {
			return user, true
		}
	}

	return nil, false
}

func (repo *UserFileRepository) GetUserByTgId(ctx context.Context, userTgId int64) (*domain.User, bool) {
	select {
	case <-ctx.Done():
		return nil, false
	default:
	}

	for _, user := range repo.items {
		if user.TelegramID == userTgId {
			return user, true
		}
	}

	return nil, false
}

func (repo *UserFileRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	user.ID = repo.nextID
	repo.nextID++

	repo.items[user.ID] = user

	err := repo.storage.Save(repo.items)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserFileRepository) Update(ctx context.Context, user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	prev := repo.items[user.ID]
	repo.items[user.ID] = user

	err := repo.storage.Save(repo.items)
	if err != nil {
		repo.items[user.ID] = prev
		return err
	}

	return nil
}

func (repo *UserFileRepository) Delete(ctx context.Context, user *domain.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	prev := user
	delete(repo.items, user.ID)

	err := repo.storage.Save(repo.items)
	if err != nil {
		repo.items[user.ID] = prev
		return err
	}

	return nil
}
