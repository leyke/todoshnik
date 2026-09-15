package auth_test

import (
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

	usererrors "todoshnik/internal/domains/user/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_SignIn(t *testing.T) {
	testUserID := 10
	testLogin := "user1"
	testPassword := "password"
	testDevice := token.DeviceTypeApi
	testRawToken := "raw-token"

	hasherError := errors.New("Hasher Error")
	errDB := errors.New("DB error")

	testUser := &user.User{
		ID:           testUserID,
		Login:        testLogin,
		PasswordHash: "hashed-pwd",
		TelegramID:   123545,
		Name:         "Имя",
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
			body:       `{"login":"` + testLogin + `","password":"` + testPassword + `"}`,
			wantStatus: http.StatusUnauthorized,
			setup: func(
				userService *mocks.UserServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
				userService.EXPECT().
					GetByLogin(mock.Anything, testLogin).
					Return(nil, errDB)
			},
		},
		{
			name:       "Успех",
			body:       `{"login":"` + testLogin + `","password":"` + testPassword + `"}`,
			wantStatus: http.StatusOK,
			want: &auth.AuthResponseDto{
				UserID:      testUserID,
				AccessToken: testRawToken,
			},
			setup: func(userService *mocks.UserServiceMock, tokenService *mocks.TokenServiceMock) {
				userService.EXPECT().GetByLogin(mock.Anything, mock.MatchedBy(func(login string) bool {
					return login == testLogin
				})).Return(testUser, nil)

				userService.EXPECT().ValidatePassword(testUser.PasswordHash, mock.MatchedBy(func(pwd string) bool {
					return pwd == testPassword
				})).Return(true, nil)

				tokenService.EXPECT().Add(mock.Anything, testUser, mock.MatchedBy(func(device token.DeviceType) bool {
					return device == testDevice
				})).Return(testRawToken, nil)
			},
		},
		{
			name:       "Пользователь не найден",
			body:       `{"login":"` + testLogin + `","password":"` + testPassword + `"}`,
			wantStatus: http.StatusUnauthorized,
			setup: func(userService *mocks.UserServiceMock, tokenService *mocks.TokenServiceMock) {
				userService.EXPECT().GetByLogin(mock.Anything, mock.MatchedBy(func(login string) bool {
					return login == testLogin
				})).Return(nil, usererrors.ErrNotFound)
			},
		},
		{
			name:       "Неверный пароль",
			body:       `{"login":"` + testLogin + `","password":"` + testPassword + `"}`,
			wantStatus: http.StatusUnauthorized,
			setup: func(userService *mocks.UserServiceMock, tokenService *mocks.TokenServiceMock) {
				userService.EXPECT().GetByLogin(mock.Anything, mock.MatchedBy(func(login string) bool {
					return login == testLogin
				})).Return(testUser, nil)

				userService.EXPECT().ValidatePassword(testUser.PasswordHash, mock.MatchedBy(func(pwd string) bool {
					return pwd == testPassword
				})).Return(false, nil)
			},
		},
		{
			name:       "Ошибка валидации пароля",
			body:       `{"login":"` + testLogin + `","password":"` + testPassword + `"}`,
			wantStatus: http.StatusUnauthorized,
			setup: func(userService *mocks.UserServiceMock, tokenService *mocks.TokenServiceMock) {
				userService.EXPECT().GetByLogin(mock.Anything, mock.MatchedBy(func(login string) bool {
					return login == testLogin
				})).Return(testUser, nil)

				userService.EXPECT().ValidatePassword(testUser.PasswordHash, mock.MatchedBy(func(pwd string) bool {
					return pwd == testPassword
				})).Return(false, hasherError)
			},
		},
		{
			name:       "Ошибка создания токена",
			body:       `{"login":"` + testLogin + `","password":"` + testPassword + `"}`,
			wantStatus: http.StatusInternalServerError,
			setup: func(userService *mocks.UserServiceMock, tokenService *mocks.TokenServiceMock) {
				userService.EXPECT().GetByLogin(mock.Anything, mock.MatchedBy(func(login string) bool {
					return login == testLogin
				})).Return(testUser, nil)

				userService.EXPECT().ValidatePassword(testUser.PasswordHash, mock.MatchedBy(func(pwd string) bool {
					return pwd == testPassword
				})).Return(true, nil)

				tokenService.EXPECT().Add(mock.Anything, testUser, testDevice).
					Return("", errDB)
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
				"/sign-in",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.SignIn(rec, req)

			if tt.want != nil {
				var got *auth.AuthResponseDto
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

				require.Equal(t, tt.want, got)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
