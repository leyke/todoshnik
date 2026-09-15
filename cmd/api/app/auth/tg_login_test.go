package auth_test

import (
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

	usererrors "todoshnik/internal/domains/user/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_TgLogin(t *testing.T) {
	testUserID := 10
	testUserTelegramID := int64(123545)

	testName := "Имя"
	testLogin := "user1"
	testDevice := token.DeviceTypeBot
	testRawToken := "raw-token"

	errDB := errors.New("DB error")

	testUser := &user.User{
		ID:           testUserID,
		Login:        testLogin,
		PasswordHash: "hashed-pwd",
		TelegramID:   testUserTelegramID,
		Name:         testName,
	}

	tests := []struct {
		name       string
		body       string
		setup      func(userService *mocks.UserServiceMock, tokenService *mocks.TokenServiceMock)
		wantStatus int
		want       *auth.AuthResponseDto
	}{
		{
			name: "невалидный json",
			body: `{"login":"`,
			setup: func(userService *mocks.UserServiceMock, tokenService *mocks.TokenServiceMock) {
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Ошибка поиска пользователя",
			body:       `{"tg_user_id":` + strconv.FormatInt(testUserTelegramID, 10) + `,"name":"` + testName + `"}`,
			wantStatus: http.StatusInternalServerError,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
				userService.EXPECT().
					GetByTgId(mock.Anything, testUserTelegramID).
					Return(nil, errDB)
			},
		},
		{
			name:       "Пользователь не найден",
			body:       `{"tg_user_id":` + strconv.FormatInt(testUserTelegramID, 10) + `,"name":"` + testName + `"}`,
			wantStatus: http.StatusNotFound,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
				userService.EXPECT().
					GetByTgId(mock.Anything, testUserTelegramID).
					Return(nil, usererrors.ErrNotFound)
			},
		},
		{
			name:       "Успех",
			body:       `{"tg_user_id":` + strconv.FormatInt(testUserTelegramID, 10) + `,"name":"` + testName + `"}`,
			wantStatus: http.StatusOK,
			want: &auth.AuthResponseDto{
				UserID:      testUserID,
				AccessToken: testRawToken,
			},
			setup: func(userService *mocks.UserServiceMock, tokenService *mocks.TokenServiceMock) {
				userService.EXPECT().
					GetByTgId(mock.Anything, testUserTelegramID).
					Return(testUser, nil)

				tokenService.EXPECT().
					Add(mock.Anything, testUser, testDevice).
					Return(testRawToken, nil)
			},
		},
		{
			name:       "Ошибка создания токена",
			body:       `{"tg_user_id":` + strconv.FormatInt(testUserTelegramID, 10) + `,"name":"` + testName + `"}`,
			wantStatus: http.StatusInternalServerError,
			setup: func(userService *mocks.UserServiceMock, tokenService *mocks.TokenServiceMock) {
				userService.EXPECT().GetByTgId(mock.Anything, mock.MatchedBy(func(tgUserID int64) bool {
					return tgUserID == testUserTelegramID
				})).Return(testUser, nil)

				tokenService.EXPECT().Add(mock.Anything, testUser, mock.MatchedBy(func(device token.DeviceType) bool {
					return device == testDevice
				})).Return("", errDB)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSvc := mocks.NewUserServiceMock(t)
			tokenSvc := mocks.NewTokenServiceMock(t)
			logger := log.New(io.Discard, "", 0)
			transactor := mocks.NewTransactorMock(t)

			tt.setup(userSvc, tokenSvc)

			handler := auth.NewHandler(userSvc, tokenSvc, logger, transactor)

			req := httptest.NewRequest(
				http.MethodPost,
				"/auth/tg/login",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.TgLogin(rec, req)

			if tt.want != nil {
				var got *auth.AuthResponseDto
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

				require.Equal(t, tt.want, got)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
