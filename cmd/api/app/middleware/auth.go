package middleware

import (
	"context"
	"net/http"
	"strings"

	authcontext "todoshnik/internal/infrastructure/context_manager/auth"
)

const (
	headerServiceToken string = "X-Bot-Service-Token"
	headerAuth         string = "Authorization"
)

type Authorizer interface {
	Authorize(ctx context.Context, token string) (int, error)
}

func Auth(authorizer Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			tokenString := tokenFromHeader(r)

			userID, err := authorizer.Authorize(r.Context(), tokenString)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := authcontext.SetUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(hfn)
	}
}

func tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get(headerAuth)
	parts := strings.Fields(headerValue)

	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return parts[1]
}

func BotAuth(botServiceToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(headerServiceToken)

			if ok := validateServiceToken(token, botServiceToken); !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hfn)
	}
}

func validateServiceToken(actual string, expected string) bool {
	return actual != "" && actual == expected
}
