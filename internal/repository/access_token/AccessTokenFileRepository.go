package accesstoken

import (
	"sync"
	"time"
	"todoshnik/internal/domain"
	"todoshnik/internal/storage"
)

type AccessTokenFileRepository struct {
	mu      sync.RWMutex
	storage storage.FileStorage[domain.Token]
	items   map[int]*domain.Token
	nextID  int
}

func NewAccessTokenFileRepository(storage storage.FileStorage[domain.Token]) (*AccessTokenFileRepository, error) {
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

	return &AccessTokenFileRepository{
		storage: storage,
		items:   items,
		nextID:  maxID + 1,
	}, nil
}

func (repo *AccessTokenFileRepository) GetAllByUserID(userID int) []*domain.Token {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	result := make([]*domain.Token, 0, len(repo.items))
	for _, token := range repo.items {
		if token.UserID == userID {
			result = append(result, token)
		}
	}
	return result
}

func (repo *AccessTokenFileRepository) GetUserIDByToken(hash string) int {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, item := range repo.items {
		if item.Hash == hash {
			return item.UserID
		}
	}
	return 0
}

func (repo *AccessTokenFileRepository) GetExpiredTokens() []*domain.Token {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	localTime := time.Now().Unix()
	result := make([]*domain.Token, 0, len(repo.items))
	for _, token := range repo.items {
		if token.ExpiresAt < localTime {
			result = append(result, token)
		}
	}
	return result
}

func (repo *AccessTokenFileRepository) Create(token *domain.Token) (*domain.Token, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	token.ID = repo.nextID
	repo.nextID++

	repo.items[token.ID] = token

	err := repo.storage.Save(repo.items)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (repo *AccessTokenFileRepository) Delete(token *domain.Token) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	prev := token
	delete(repo.items, token.ID)

	err := repo.storage.Save(repo.items)
	if err != nil {
		repo.items[token.ID] = prev
		return err
	}

	return nil
}
