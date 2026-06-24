package user

import (
	"context"
	"sort"
	"sync"

	userdomain "todoshnik/internal/domains/user"
	filestorage "todoshnik/internal/infrastructure/filestorage"
)

type FileRepository struct {
	mu      sync.RWMutex
	storage filestorage.FileStorage[userdomain.User]
	items   map[int]*userdomain.User
	nextID  int
}

func NewFileRepository(storage filestorage.FileStorage[userdomain.User]) (*FileRepository, error) {
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

	return &FileRepository{
		storage: storage,
		items:   items,
		nextID:  maxID + 1,
	}, nil
}

func (repo *FileRepository) List(ctx context.Context) []*userdomain.User {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

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

	result := make([]*userdomain.User, 0, len(keys))
	for _, k := range keys {
		result = append(result, repo.items[k])
	}
	return result
}

func (repo *FileRepository) GetByID(ctx context.Context, id int) (*userdomain.User, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, false
	default:
	}

	user, ok := repo.items[id]
	if !ok {
		return nil, false
	}

	copy := *user
	return &copy, ok
}

func (repo *FileRepository) GetByLogin(ctx context.Context, login string) (*userdomain.User, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, false
	default:
	}

	for _, user := range repo.items {
		if user.Login == login {
			copy := *user
			return &copy, true
		}
	}

	return nil, false
}

func (repo *FileRepository) GetByTgId(ctx context.Context, userTgId int64) (*userdomain.User, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, false
	default:
	}

	for _, user := range repo.items {
		if user.TelegramID == userTgId {
			copy := *user
			return &copy, true
		}
	}

	return nil, false
}

func (repo *FileRepository) Create(ctx context.Context, user *userdomain.User) (*userdomain.User, error) {
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

func (repo *FileRepository) Update(ctx context.Context, user *userdomain.User) error {
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

func (repo *FileRepository) Delete(ctx context.Context, user *userdomain.User) error {
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
