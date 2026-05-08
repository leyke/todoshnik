package user

import (
	"context"
	"strconv"
	"time"
	"todoshnik/internal/bot/tg"

	"github.com/redis/go-redis/v9"
)

type UserState struct {
	tg.Command
	tg.State
}

type StateStorage struct {
	rdb *redis.Client
}

func NewStateStorage(rdb *redis.Client) *StateStorage {
	return &StateStorage{
		rdb: rdb,
	}
}

func (ss *StateStorage) Set(ctx context.Context, userID int64, command tg.Command, state tg.State) error {
	key := getKey(userID)

	ss.rdb.HSet(ctx, key,
		"command", command,
		"state", state,
	)

	// запоминаем команду на час
	err := ss.rdb.Expire(ctx, key, 1*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ss *StateStorage) Get(ctx context.Context, userID int64) (*UserState, bool) {
	data, err := ss.rdb.HGetAll(ctx, getKey(userID)).Result()
	if err != nil {
		return nil, false
	}

	// Проверяем, существует ли пользователь
	if len(data) == 0 {
		return nil, false
	}

	// Создаем структуру и заполняем
	us := &UserState{}
	if command, ok := data["command"]; ok {
		us.Command = tg.Command(command)
		return nil, false
	}

	if state, ok := data["state"]; ok {
		us.State = tg.State(state)
	}

	return us, true
}

func getKey(userID int64) string {
	return "user:" + strconv.FormatInt(userID, 10) + ":tg-command-state"
}
