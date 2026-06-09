package session

import (
	"context"
	"fmt"
	"time"

	"todoshnik/internal/bot/tg"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	rdb *redis.Client
}

func NewStorage(rdb *redis.Client) *Storage {
	return &Storage{
		rdb: rdb,
	}
}

// Аналогичная претензия– надо вынести инфраструктурную обертку над редисом за интерфейс и на инфраструктурный слой

func (ss *Storage) Set(ctx context.Context, userID int64, command tg.Command, state tg.State) error {
	key := getKey(userID)

	err := ss.rdb.HSet(ctx, key,
		"command", string(command),
		"state", string(state),
	).Err()
	if err != nil {
		return err
	}
	// запоминаем команду на час
	err = ss.rdb.Expire(ctx, key, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ss *Storage) Get(ctx context.Context, userID int64) (*Session, bool) {
	key := getKey(userID)
	data, err := ss.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	if len(data) == 0 {
		return nil, false
	}

	us := &Session{
		State: tg.StateIdle,
	}
	command, ok := data["command"]
	if !ok {
		return nil, false
	}
	us.Command = tg.Command(command)

	if state, ok := data["state"]; ok {
		us.State = tg.State(state)
	}

	return us, true
}

// формат ключа user:{userID}, чтобы можно было проще удалять по id юзера через регулярку rk.xf, если понадобится
func getKey(userID int64) string {
	return fmt.Sprintf("user:%d:tg-command-state", userID)
}
