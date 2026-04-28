package service

import (
	"os"
	"strconv"
	"time"

	"todoshnik/internal/auth"
	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	repository "todoshnik/internal/repository/access_token"
	"todoshnik/internal/storage"
)

const (
	DefaultTokenTtl int = 14
)

type AccessTokenService struct {
	repo repository.AccessTokenRepositoryInterface
}

func NewAccessTokenService() (*AccessTokenService, error) {
	storagePath := os.Getenv("tmp_dir") + "/tokens.json"
	storage := storage.NewFileStorage[domain.Token](storagePath)

	repo, err := repository.NewAccessTokenFileRepository(storage)
	if err != nil {
		return nil, err
	}

	s := &AccessTokenService{
		repo: repo,
	}

	return s, nil
}

func (s *AccessTokenService) AddToken(user domain.User) (*domain.Token, error) {
	token, tokenError := auth.GenerateToken()
	if tokenError != nil {
		return nil, tokenError
	}

	hash := auth.HashToken(token, os.Getenv("salt"))

	tokenTtl, tokenError := strconv.Atoi(os.Getenv("token_ttl_days"))
	if tokenError != nil {
		tokenTtl = DefaultTokenTtl
	}
	localTime := time.Now().AddDate(0, 0, tokenTtl).Unix()

	newToken := &domain.Token{
		UserID:    user.ID,
		Hash:      hash,
		ExpiresAt: localTime,
	}

	newToken, err := s.repo.Create(newToken)
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func (s *AccessTokenService) GetUserID(token string) (int, error) {
	hash := auth.HashToken(token, os.Getenv("salt"))
	userID := s.repo.GetUserIDByToken(hash)

	if userID == 0 {
		return 0, apperrors.ErrNotFound
	}

	return userID, nil
}

func (s *AccessTokenService) ClearExpiredTokens() int {
	tokens := s.repo.GetExpiredTokens()
	counter := 0
	for _, token := range tokens {
		err := s.repo.Delete(token)
		if err != nil {
			counter++
		}
	}

	return counter
}
