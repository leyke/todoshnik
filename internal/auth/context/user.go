package context

import (
	"context"

	"todoshnik/internal/user"
)

func SetUser(ctx context.Context, user *user.User) context.Context {
	ctx = context.WithValue(ctx, UserID, user.ID)
	ctx = context.WithValue(ctx, User, user)

	return ctx
}

//  1. Твой пакет контекст зависит от типа юзер, что печально. В целом это должно быть либо частью пакета user либо
//     подчиненным для пакета user пакетом: user/metadata
func GetUser(ctx context.Context) (*user.User, bool) {
	user, ok := ctx.Value(User).(*user.User)
	return user, ok
}

func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(UserID).(int)
	return id, ok
}
