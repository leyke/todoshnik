package service

import (
	"os"
	"strconv"
	"time"

	"todoshnik/internal/auth"
	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	repository "todoshnik/internal/repository/access_token"
)

const (
	DefaultTokenTtl int = 14 // количество дней жизни токена
)

type AccessTokenService struct {
	repo repository.AccessTokenRepositoryInterface
}

func NewAccessTokenService(repo repository.AccessTokenRepositoryInterface) *AccessTokenService {
	return &AccessTokenService{
		repo: repo,
	}
}

func (s *AccessTokenService) AddToken(user *domain.User, device domain.DeviceType) (*domain.Token, error) {
	token, tokenError := auth.GenerateToken()
	if tokenError != nil {
		return nil, tokenError
	}

	hash := auth.HashToken(token, os.Getenv("SALT"))

	tokenTtl, tokenError := strconv.Atoi(os.Getenv("TOKEN_TTL_DAYS"))
	if tokenError != nil {
		tokenTtl = DefaultTokenTtl
	}
	localTime := time.Now().AddDate(0, 0, tokenTtl).Unix()

	newToken := &domain.Token{
		UserID:    user.ID,
		Hash:      hash,
		ExpiresAt: localTime,
		Device:    device,
	}

	newToken, err := s.repo.Create(newToken)
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func (s *AccessTokenService) GetUserID(token string) (int, error) {
	userID := s.repo.GetUserIDByToken(token)

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
