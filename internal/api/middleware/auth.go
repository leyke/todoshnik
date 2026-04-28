package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"todoshnik/internal/api/contextkeys"
	"todoshnik/internal/auth"
	"todoshnik/internal/service"
)

func Auth(ats *service.AccessTokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, ok := extractTokenFromHeader(r)
			if !ok {
				http.Error(w, "Unauthorized: не передан токен", http.StatusUnauthorized)
				return
			}

			hash := auth.HashToken(tokenString, os.Getenv("salt"))
			userID, err := ats.GetUserID(hash)
			if err != nil {
				http.Error(w, "Unauthorized: не найден пользователь", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromHeader(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}

	return token, true
}
