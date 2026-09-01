package bot

import (
	"context"
	"fmt"
	"time"

	"todoshnik/cmd/tg_bot/app/bot/command"

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

func (s *Storage) StartCommand(ctx context.Context, userID int64, c command.Name) error {
	return s.set(ctx, userID, c, command.StateWait)
}

func (s *Storage) FinishCommand(ctx context.Context, userID int64, c command.Name) error {
	return s.set(ctx, userID, c, command.StateComplete)
}

func (s *Storage) GetLastCommand(ctx context.Context, userID int64) (*command.CommandDto, bool) {
	return s.get(ctx, userID)
}

func (s *Storage) set(ctx context.Context, userID int64, c command.Name, state command.State) error {
	key := getKey(userID)

	err := s.rdb.HSet(ctx, key,
		"command", string(c),
		"state", string(state),
	).Err()
	if err != nil {
		return err
	}
	// запоминаем команду на час
	err = s.rdb.Expire(ctx, key, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) get(ctx context.Context, userID int64) (*command.CommandDto, bool) {
	key := getKey(userID)
	data, err := s.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	if len(data) == 0 {
		return nil, false
	}

	dto := &command.CommandDto{
		State: command.StateIdle,
	}
	c, ok := data["command"]
	if !ok {
		return nil, false
	}
	dto.Name = command.Name(c)

	if state, ok := data["state"]; ok {
		dto.State = command.State(state)
	}

	return dto, true
}

// формат ключа user:{userID}, чтобы можно было проще удалять по id юзера через регулярку rk.xf, если понадобится
func getKey(userID int64) string {
	return fmt.Sprintf("user:%d:tg-command-state", userID)
}
