package redis

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func NewClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
