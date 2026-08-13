package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"todoshnik/cmd/api/app/auth"
	"todoshnik/cmd/api/app/auth/mocks"
	"todoshnik/internal/domains/token"
	"todoshnik/internal/domains/user"
	"todoshnik/internal/infrastructure/validation"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_TgAutoReg(t *testing.T) {
	testUserID := 10
	testTelegramID := int64(123545)
	testName := "Имя"
	testRawToken := "raw-token"

	errDB := errors.New("DB error")

	testUser := &user.User{
		ID:         testUserID,
		TelegramID: testTelegramID,
		Name:       testName,
	}

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
			name:       "невалидный json",
			body:       `{"tg_user_id":"`,
			wantStatus: http.StatusBadRequest,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
				transactor *mocks.TransactorMock,
			) {
			},
		},
		{
			name: "ошибка создания пользователя",
			body: `{"tg_user_id":` + strconv.FormatInt(testTelegramID, 10) +
				`,"name":"` + testName + `"}`,
			wantStatus: http.StatusInternalServerError,
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
					AddFromTg(mock.Anything, testName, testTelegramID).
					Return(nil, errDB)
			},
		},
		{
			name: "невалидные данные пользователя",
			body: `{"tg_user_id":` + strconv.FormatInt(testTelegramID, 10) +
				`,"name":"` + testName + `"}`,
			wantStatus: http.StatusBadRequest,
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
					AddFromTg(mock.Anything, testName, testTelegramID).
					Return(nil, validation.ErrNotValidate)
			},
		},
		{
			name: "ошибка создания токена",
			body: `{"tg_user_id":` + strconv.FormatInt(testTelegramID, 10) +
				`,"name":"` + testName + `"}`,
			wantStatus: http.StatusInternalServerError,
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
					AddFromTg(mock.Anything, testName, testTelegramID).
					Return(testUser, nil)

				tokenService.EXPECT().
					Add(mock.Anything, testUser, token.DeviceTypeBot).
					Return("", errDB)
			},
		},
		{
			name: "успех",
			body: `{"tg_user_id":` + strconv.FormatInt(testTelegramID, 10) +
				`,"name":"` + testName + `"}`,
			wantStatus: http.StatusOK,
			want: &auth.AuthResponseDto{
				UserID:      testUserID,
				AccessToken: testRawToken,
			},
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
					AddFromTg(mock.Anything, testName, testTelegramID).
					Return(testUser, nil)

				tokenService.EXPECT().
					Add(mock.Anything, testUser, token.DeviceTypeBot).
					Return(testRawToken, nil)
			},
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
				"/auth/tg/auto-reg",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.TgAutoReg(rec, req)

			if tt.want != nil {
				var got *auth.AuthResponseDto

				require.NoError(
					t,
					json.Unmarshal(rec.Body.Bytes(), &got),
				)

				require.Equal(t, tt.want, got)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
