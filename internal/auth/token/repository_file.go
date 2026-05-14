package token

import (
	"context"
	"sync"
	"time"
	"todoshnik/internal/storage"
)

type FileRepository struct {
	mu      sync.RWMutex
	storage storage.FileStorage[Token]
	items   map[int]*Token
	nextID  int
}

func NewFileRepository(storage storage.FileStorage[Token]) (*FileRepository, error) {
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

func (repo *FileRepository) GetAllByUserID(ctx context.Context, userID int) []*Token {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	result := make([]*Token, 0, len(repo.items))
	for _, token := range repo.items {
		if token.UserID == userID {
			result = append(result, token)
		}
	}
	return result
}

func (repo *FileRepository) GetByHash(ctx context.Context, hash string) *Token {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	for _, item := range repo.items {
		if item.Hash == hash {
			copy := *item
			return &copy
		}
	}
	return nil
}

func (repo *FileRepository) GetExpiredTokens(ctx context.Context) []*Token {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	localTime := time.Now().Unix()
	result := make([]*Token, 0, len(repo.items))
	for _, token := range repo.items {
		if token.ExpiresAt < localTime {
			result = append(result, token)
		}
	}
	return result
}

func (repo *FileRepository) Create(ctx context.Context, token *Token) (*Token, error) {
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

func (repo *FileRepository) Delete(ctx context.Context, token *Token) error {
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
