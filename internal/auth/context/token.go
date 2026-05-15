package context

import (
	"context"
)

func SetToken(ctx context.Context, token string) context.Context {
	ctx = context.WithValue(ctx, Token, token)
	return ctx
}

func GetToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(Token).(string)
	return token, ok
}
