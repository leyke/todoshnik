package token

import "context"

type Storage interface {
	Get(ctx context.Context, userID int64) (string, error)
	Set(ctx context.Context, userID int64, token string) error
}
