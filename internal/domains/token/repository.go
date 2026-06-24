package token

import (
	"context"
)

type Repository interface {
	GetAllByUserID(ctx context.Context, id int) ([]*Token, error)
	GetByHash(ctx context.Context, hash string) (*Token, error)
	GetExpiredTokens(ctx context.Context) ([]*Token, error)

	Create(ctx context.Context, token *Token) (*Token, error)
	Delete(ctx context.Context, token *Token) error
}
