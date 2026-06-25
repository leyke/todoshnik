package context

import (
	"context"
)

func SetToken(ctx context.Context, token string) context.Context {
	ctx = context.WithValue(ctx, accessTokenKey, token)
	return ctx
}

func GetToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(accessTokenKey).(string)
	return token, ok
}
