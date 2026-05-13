package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	ah "todoshnik/internal/api/auth"
	"todoshnik/internal/api/contextkeys"
)

func Auth(ah *ah.AuthHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, ok := extractTokenFromHeader(r)
			if !ok {
				http.Error(w, "Unauthorized: не передан токен", http.StatusUnauthorized)
				return
			}

			user, err := ah.ValidateToken(r.Context(), tokenString)
			if err != nil {
				fmt.Printf("Ошибка валидации токена: %v\n", err)
				http.Error(w, "Unauthorized: не найден пользователь", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, user.ID)
			ctx = context.WithValue(ctx, contextkeys.UserKey, user)
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

func BotAuth(ah *ah.AuthHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("X-Bot-Service-Token")

			token := strings.TrimSpace(authHeader)
			if token == "" {
				http.Error(w, "Unauthorized: не передан токен", http.StatusUnauthorized)
				return
			}

			if token != os.Getenv("BOT_SERVICE_TOKEN") {
				fmt.Printf("Неверный сервисный токен бота")
				http.Error(w, "Unauthorized: сервисный токен не верен", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
