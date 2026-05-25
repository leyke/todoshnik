package middleware

import (
	"context"
	"strings"
	authcontext "todoshnik/internal/auth/context"
	"todoshnik/internal/auth/token"
	"todoshnik/internal/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(userService *user.Service, tokenService *token.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(
				codes.Unauthenticated,
				"metadata required",
			)
		}

		values := md.Get("authorization")

		if len(values) == 0 {
			return nil, status.Error(
				codes.Unauthenticated,
				"token required",
			)
		}

		rawToken := strings.TrimPrefix(
			values[0],
			"Bearer ",
		)

		token, err := tokenService.Get(ctx, rawToken)
		if err != nil {
			return nil, err
		}

		user, err := userService.Get(ctx, token.UserID, "")
		if err != nil {
			return nil, err
		}

		ctx = authcontext.SetUser(ctx, user)

		return handler(ctx, req)
	}
}
