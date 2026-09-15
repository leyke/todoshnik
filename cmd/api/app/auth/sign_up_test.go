package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todoshnik/cmd/api/app/auth"
	"todoshnik/cmd/api/app/auth/mocks"
	"todoshnik/internal/domains/token"
	"todoshnik/internal/domains/user"
	"todoshnik/internal/infrastructure/validation"

	usererrors "todoshnik/internal/domains/user/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_SignUp(t *testing.T) {
	testUserID := 10
	testName := "Имя"
	testLogin := "user1"
	testPassword := "password"
	testRawToken := "raw-token"
	testDevice := token.DeviceTypeApi

	errDB := errors.New("DB error")

	testUser := &user.User{
		ID:           testUserID,
		Login:        testLogin,
		PasswordHash: "hashed-pwd",
		TelegramID:   123545,
		Name:         testName,
	}

	testBody := `{"name":"` + testName + `","login":"` + testLogin + `","password":"` + testPassword + `"}`

	tests := []struct {
		name  string
		body  string
		setup func(
			userService *mocks.UserServiceMock,
			tokenService *mocks.TokenServiceMock,
			transactor *mocks.TransactorMock,
		)
		wantStatus int
		want       *auth.AuthResponseDto
	}{
		{
			name: "невалидный json",
			body: `{"login":"`,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
				transactor *mocks.TransactorMock,
			) {
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Успех",
			body: testBody,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
				transactor *mocks.TransactorMock,
			) {
				transactor.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						fn func(context.Context) error,
					) error {
						return fn(ctx)
					})

				userService.EXPECT().
					Add(mock.Anything, testName, testLogin, testPassword).
					Return(testUser, nil)

				tokenService.EXPECT().
					Add(mock.Anything, testUser, testDevice).
					Return(testRawToken, nil)
			},
			wantStatus: http.StatusOK,
			want: &auth.AuthResponseDto{
				UserID:      testUserID,
				AccessToken: testRawToken,
			},
		},
		{
			name: "Пользователь уже существует",
			body: testBody,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
				transactor *mocks.TransactorMock,
			) {
				transactor.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						fn func(context.Context) error,
					) error {
						return fn(ctx)
					})

				userService.EXPECT().
					Add(mock.Anything, testName, testLogin, testPassword).
					Return(nil, usererrors.ErrConflict)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "Ошибка валидации пользователя",
			body: testBody,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
				transactor *mocks.TransactorMock,
			) {
				transactor.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						fn func(context.Context) error,
					) error {
						return fn(ctx)
					})

				userService.EXPECT().
					Add(mock.Anything, testName, testLogin, testPassword).
					Return(nil, validation.NewValidationErrorFromValidator(nil))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Ошибка создания пользователя",
			body: testBody,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
				transactor *mocks.TransactorMock,
			) {
				transactor.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						fn func(context.Context) error,
					) error {
						return fn(ctx)
					})

				userService.EXPECT().
					Add(mock.Anything, testName, testLogin, testPassword).
					Return(nil, errDB)
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Ошибка создания токена",
			body: testBody,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
				transactor *mocks.TransactorMock,
			) {
				transactor.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(
						ctx context.Context,
						fn func(context.Context) error,
					) error {
						return fn(ctx)
					})

				userService.EXPECT().
					Add(mock.Anything, testName, testLogin, testPassword).
					Return(testUser, nil)

				tokenService.EXPECT().
					Add(mock.Anything, testUser, testDevice).
					Return("", errDB)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSvc := mocks.NewUserServiceMock(t)
			tokenSvc := mocks.NewTokenServiceMock(t)
			transactor := mocks.NewTransactorMock(t)
			logger := log.New(io.Discard, "", 0)

			tt.setup(userSvc, tokenSvc, transactor)

			handler := auth.NewHandler(
				userSvc,
				tokenSvc,
				logger,
				transactor,
			)

			req := httptest.NewRequest(
				http.MethodPost,
				"/sign-up",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.SignUp(rec, req)

			if tt.want != nil {
				var got *auth.AuthResponseDto

				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

				require.Equal(t, tt.want, got)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
