package token

import (
	"context"
	"os"
	"strconv"
	"time"

	"todoshnik/internal/auth"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/user"
)

const (
	DefaultTokenTtl int = 14 // количество дней жизни токена
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(ctx context.Context, user *user.User, device auth.DeviceType) (string, error) {
	token, tokenError := GenerateToken()
	if tokenError != nil {
		return "", tokenError
	}

	hash := HashToken(token, os.Getenv("SALT"))

	tokenTtl, tokenError := strconv.Atoi(os.Getenv("TOKEN_TTL_DAYS"))
	if tokenError != nil {
		tokenTtl = DefaultTokenTtl
	}
	localTime := time.Now().AddDate(0, 0, tokenTtl).Unix()

	newToken := &Token{
		UserID:    user.ID,
		Hash:      hash,
		ExpiresAt: localTime,
		Device:    device,
	}

	newToken, err := s.repo.Create(ctx, newToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) Get(ctx context.Context, rawToken string) (*Token, error) {
	hash := HashToken(rawToken, os.Getenv("SALT"))

	token := s.repo.GetByHash(ctx, hash)

	if token == nil {
		return nil, apperrors.ErrNotFound
	}

	return token, nil
}

func (s *Service) ClearExpiredTokens(ctx context.Context) int {
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
