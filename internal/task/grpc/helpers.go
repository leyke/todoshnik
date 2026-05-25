package grpc

import (
	"context"
	authcontext "todoshnik/internal/auth/context"
	"todoshnik/internal/identity"

	"google.golang.org/grpc/metadata"
)

func getScope(userID int) identity.AccessScope {
	return identity.AccessScope{
		IsAdmin: false,
		UserID:  userID,
	}
}

func authContext(
	ctx context.Context,
) context.Context {
	token, ok := authcontext.GetToken(ctx)
	if !ok {
		return ctx
	}

	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	return metadata.NewOutgoingContext(ctx, md)
}
