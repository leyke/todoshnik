package accesstoken

import (
	"context"
	"todoshnik/internal/domain"
)

type AccessTokenRepositoryInterface interface {
	GetAllByUserID(ctx context.Context, id int) []*domain.Token
	GetUserIDByToken(ctx context.Context, token string) int
	GetExpiredTokens(ctx context.Context) []*domain.Token

	Create(ctx context.Context, token *domain.Token) (*domain.Token, error)
	Delete(ctx context.Context, token *domain.Token) error
}
