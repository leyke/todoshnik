package token

import (
	"context"
	"errors"

	userdomain "todoshnik/internal/domains/user"
)

type Service struct {
	repo           Repository
	tokenGenerator TokenGenerator
	clock          Clock
	cfg            Config
}

func NewService(repo Repository, tokenGenerator TokenGenerator, clock Clock, cfg Config) *Service {
	return &Service{
		repo:           repo,
		tokenGenerator: tokenGenerator,
		clock:          clock,
		cfg:            cfg,
	}
}

func (s *Service) Add(ctx context.Context, user *userdomain.User, device DeviceType) (string, error) {
	token, tokenError := s.tokenGenerator.Generate()
	if tokenError != nil {
		return "", tokenError
	}

	hash := s.tokenGenerator.Hash(token)
	localTime := s.clock.Now().Add(s.cfg.Ttl)

	_, err := s.repo.Create(ctx, &Token{
		UserID:    user.ID,
		Hash:      hash,
		ExpiresAt: localTime,
		Device:    device,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Get(ctx context.Context, rawToken string) (*Token, error) {
	hash := s.tokenGenerator.Hash(rawToken)

	token, err := s.repo.GetByHash(ctx, hash)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) ClearExpiredTokens(ctx context.Context) (int, error) {
	tokens, err := s.repo.GetExpiredTokens(ctx, s.clock.Now())
	if err != nil {
		return 0, err
	}
	deletedCount := 0
	errs := make([]error, 0, len(tokens))
	for _, token := range tokens {
		err := s.repo.Delete(ctx, token)
		if err != nil {
			errs = append(errs, err)
		} else {
			deletedCount++
		}
	}

	return deletedCount, errors.Join(errs...)
}
