package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"todoshnik/cmd/api/app/middleware"
	"todoshnik/cmd/api/app/middleware/mocks"

	authcontext "todoshnik/internal/infrastructure/context_manager/auth"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Auth(t *testing.T) {
	tests := []struct {
		name       string
		header     string
		setup      func(authorizer *mocks.AuthorizerMock)
		wantStatus int
		wantUserID int
	}{
		{
			name:   "успех",
			header: "Bearer raw-token",
			setup: func(authorizer *mocks.AuthorizerMock) {
				authorizer.
					On("Authorize", mock.Anything, "raw-token").
					Return(10, nil)
			},
			wantStatus: http.StatusOK,
			wantUserID: 10,
		},
		{
			name:   "ошибка авторизации",
			header: "Bearer raw-token",
			setup: func(authorizer *mocks.AuthorizerMock) {
				authorizer.
					On("Authorize", mock.Anything, "raw-token").
					Return(0, errors.New("invalid token"))
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "нет заголовка",
			header: "",
			setup: func(authorizer *mocks.AuthorizerMock) {
				authorizer.
					On("Authorize", mock.Anything, "").
					Return(0, errors.New("empty token"))
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "не Bearer",
			header: "Basic raw-token",
			setup: func(authorizer *mocks.AuthorizerMock) {
				authorizer.
					On("Authorize", mock.Anything, "").
					Return(0, errors.New("empty token"))
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "неверный формат",
			header: "Bearer",
			setup: func(authorizer *mocks.AuthorizerMock) {
				authorizer.
					On("Authorize", mock.Anything, "").
					Return(0, errors.New("empty token"))
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authorizer := mocks.NewAuthorizerMock(t)

			tt.setup(authorizer)

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, ok := authcontext.GetUserID(r.Context())

				require.True(t, ok)
				require.Equal(t, tt.wantUserID, userID)

				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.Auth(authorizer)(next)

			req := httptest.NewRequest(
				http.MethodGet,
				"/tasks",
				nil,
			)

			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			authorizer.AssertExpectations(t)
		})
	}
}

func Test_BotAuth(t *testing.T) {
	const botServiceToken = "bot-service-token"

	tests := []struct {
		name       string
		header     string
		wantStatus int
	}{
		{
			name:       "успех",
			header:     botServiceToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "нет заголовка",
			header:     "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "неверный токен",
			header:     "invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.BotAuth(botServiceToken)(next)

			req := httptest.NewRequest(
				http.MethodGet,
				"/auth/tg/login",
				nil,
			)

			if tt.header != "" {
				req.Header.Set("X-Bot-Service-Token", tt.header)
			}

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
