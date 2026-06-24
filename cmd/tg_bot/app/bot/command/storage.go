package command

import "context"

type Storage interface {
	GetLastCommand(ctx context.Context, userID int64) (*CommandDto, bool)

	StartCommand(ctx context.Context, userID int64, c Name) error
	FinishCommand(ctx context.Context, userID int64, c Name) error
}
