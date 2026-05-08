package user

import (
	"context"
	"fmt"
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
	fmt.Println(key)

	err := ss.rdb.HSet(ctx, key,
		"command", string(command),
		"state", string(state),
	).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	// запоминаем команду на час
	err = ss.rdb.Expire(ctx, key, 1*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ss *StateStorage) Get(ctx context.Context, userID int64) (*UserState, bool) {
	key := getKey(userID)
	fmt.Println(key)
	data, err := ss.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, false
	}
	fmt.Println(data)
	if len(data) == 0 {
		return nil, false
	}

	us := &UserState{
		State: tg.StateIdale,
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

func getKey(userID int64) string {
	return "user:" + strconv.FormatInt(userID, 10) + ":tg-command-state"
}
