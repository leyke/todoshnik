package token

import (
	"context"
)

type Repository interface {
	GetAllByUserID(ctx context.Context, id int) []*Token
	GetByHash(ctx context.Context, hash string) *Token
	GetExpiredTokens(ctx context.Context) []*Token

	Create(ctx context.Context, token *Token) (*Token, error)
	Delete(ctx context.Context, token *Token) error
}
