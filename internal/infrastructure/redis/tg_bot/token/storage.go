package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	rdb      *redis.Client
	tokenTTL time.Duration
}

func NewStorage(rdb *redis.Client) *Storage {
	return &Storage{
		rdb:      rdb,
		tokenTTL: 4 * time.Hour,
	}
}

func (s *Storage) Set(ctx context.Context, userID int64, token string) error {
	key := getKey(userID)

	err := s.rdb.Set(
		ctx,
		key,
		token,
		s.tokenTTL,
	).Err()

	return err
}

func (s *Storage) Get(ctx context.Context, userID int64) (string, error) {
	key := getKey(userID)
	token, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	return token, err
}

// формат ключа user:{userID}, чтобы можно было проще удалять по id юзера через регулярку rk.xf, если понадобится
func getKey(userID int64) string {
	return fmt.Sprintf("user:%d:tg-api-token-key", userID)
}
