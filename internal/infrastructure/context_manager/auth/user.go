package context

import (
	"context"
)

func SetUserID(ctx context.Context, userID int) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)

	return ctx
}

func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}
