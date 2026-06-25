package token

import (
	"context"
	"errors"
	"time"

	userdomain "todoshnik/internal/domains/user"
)

type Service struct {
	repo Repository
	cfg  Config
}

func NewService(repo Repository, cfg Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *Service) Add(ctx context.Context, user *userdomain.User, device DeviceType) (string, error) {
	token, tokenError := GenerateToken()
	if tokenError != nil {
		return "", tokenError
	}

	hash := HashToken(token, s.cfg.Secret)
	localTime := time.Now().Add(s.cfg.Ttl).Unix()

	tokenObject := &Token{
		UserID:    user.ID,
		Hash:      hash,
		ExpiresAt: localTime,
		Device:    device,
	}

	_, err := s.repo.Create(ctx, tokenObject)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) Get(ctx context.Context, rawToken string) (*Token, error) {
	hash := HashToken(rawToken, s.cfg.Secret)

	token, err := s.repo.GetByHash(ctx, hash)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) ClearExpiredTokens(ctx context.Context) (int, error) {
	tokens, err := s.repo.GetExpiredTokens(ctx)
	if err != nil {
		return 0, err
	}
	counter := 0
	errs := make([]error, 0, len(tokens))
	for _, token := range tokens {
		err := s.repo.Delete(ctx, token)
		if err != nil {
			counter++
			errs = append(errs, err)
		}
	}

	return counter, errors.Join(errs...)
}
