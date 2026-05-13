package service

import (
	"context"
	"os"
	"strconv"
	"time"

	"todoshnik/internal/auth"
	"todoshnik/internal/domain"
	apperrors "todoshnik/internal/errors"
	repository "todoshnik/internal/repository/access_token"
	"todoshnik/internal/user"
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

func (s *AccessTokenService) AddToken(ctx context.Context, user *user.User, device domain.DeviceType) (*domain.Token, error) {
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

	newToken, err := s.repo.Create(ctx, newToken)
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func (s *AccessTokenService) GetUserID(ctx context.Context, token string) (int, error) {
	userID := s.repo.GetUserIDByToken(ctx, token)

	if userID == 0 {
		return 0, apperrors.ErrNotFound
	}

	return userID, nil
}

func (s *AccessTokenService) ClearExpiredTokens(ctx context.Context) int {
	tokens := s.repo.GetExpiredTokens(ctx)
	counter := 0
	for _, token := range tokens {
		err := s.repo.Delete(ctx, token)
		if err != nil {
			counter++
		}
	}

	return counter
}
