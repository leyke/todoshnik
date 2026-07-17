package token

import (
	"context"
	"sync"
	"time"

	tokendomain "todoshnik/internal/domains/token"
	filestorage "todoshnik/internal/infrastructure/filestorage"
)

type FileRepository struct {
	mu      sync.RWMutex
	storage filestorage.FileStorage[tokendomain.Token]
	items   map[int]*tokendomain.Token
	nextID  int
}

func NewFileRepository(storage filestorage.FileStorage[tokendomain.Token]) (*FileRepository, error) {
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

func (repo *FileRepository) GetAllByUserID(ctx context.Context, userID int) ([]*tokendomain.Token, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	result := make([]*tokendomain.Token, 0, len(repo.items))
	for _, token := range repo.items {
		if token.UserID == userID {
			result = append(result, token)
		}
	}
	return result, nil
}

func (repo *FileRepository) GetByHash(ctx context.Context, hash string) (*tokendomain.Token, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	for _, item := range repo.items {
		if item.Hash == hash {
			copy := *item
			return &copy, nil
		}
	}
	return nil, nil
}

func (repo *FileRepository) GetExpiredTokens(ctx context.Context) ([]*tokendomain.Token, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	localTime := time.Now()

	result := make([]*tokendomain.Token, 0, len(repo.items))
	for _, token := range repo.items {
		if token.ExpiresAt.Before(localTime) {
			result = append(result, token)
		}
	}
	return result, nil
}

func (repo *FileRepository) Create(ctx context.Context, token *tokendomain.Token) (*tokendomain.Token, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	token.ID = repo.nextID
	repo.nextID++

	repo.items[token.ID] = token

	err := repo.storage.Save(repo.items)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (repo *FileRepository) Delete(ctx context.Context, token *tokendomain.Token) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	prev := token
	delete(repo.items, token.ID)

	err := repo.storage.Save(repo.items)
	if err != nil {
		repo.items[token.ID] = prev
		return err
	}

	return nil
}
