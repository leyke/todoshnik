package redis

import (
	"context"
	"fmt"
	"time"
	"todoshnik/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.Config) (*redis.Client, error) {
	// решена ли проблема холодного старта? Как будет работать приложение/инстанс приложения если один из подов
	// перезагрузится?
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return rdb, nil
}
