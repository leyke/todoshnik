package session

import (
	"context"
	"fmt"
	"strconv"
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
	// интревал в настройки. 1 можно не писать, просто time.Hour
	err = ss.rdb.Expire(ctx, key, 1*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ss *Storage) Get(ctx context.Context, userID int64) (*Session, bool) {
	key := getKey(userID)
	data, err := ss.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		// у тебя же был где-то логер
		fmt.Println("Ошибка получения данных из Redis:", err)
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

func getKey(userID int64) string {
	// 1. если ключ будет "user:tg-command-state:%id" то будет на 1 конкатинацию меньше
	// 2. можно использовать fmt было бы читабельнее и возможно производительнее (надо замерять)
	return "user:" + strconv.FormatInt(userID, 10) + ":tg-command-state"
}
