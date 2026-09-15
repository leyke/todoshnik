package auth_test

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"
	"time"

	"todoshnik/cmd/api/app/auth"
	"todoshnik/cmd/api/app/auth/mocks"
	"todoshnik/internal/domains/token"

	apierror "todoshnik/cmd/api/app/errors"
	tokenerrors "todoshnik/internal/domains/token/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Authorize(t *testing.T) {
	testUserID := 10
	testToken := &token.Token{
		ID:        1,
		UserID:    testUserID,
		Hash:      "hash-token",
		Device:    token.DeviceTypeApi,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	errDB := errors.New("DB error")

	tests := []struct {
		name     string
		rawToken string
		setup    func(tokenService *mocks.TokenServiceMock, rt string)
		wantErr  error
		want     int
	}{
		{
			name:     "Успешно получить юзера по токену",
			rawToken: "raw-token",
			setup: func(tokenService *mocks.TokenServiceMock, rt string) {
				tokenService.EXPECT().Get(mock.Anything, mock.MatchedBy(func(rawToken string) bool {
					return rawToken == rt
				})).Return(testToken, nil)
			},
			want: testUserID,
		},
		{
			name:     "Пришел пустой токен",
			rawToken: "",
			setup: func(tokenService *mocks.TokenServiceMock, rt string) {

			},
			wantErr: apierror.ErrUnauth,
			want:    0,
		},
		{
			name:     "Нет токена в базе",
			rawToken: "rnd-token",
			setup: func(tokenService *mocks.TokenServiceMock, rt string) {
				tokenService.EXPECT().Get(mock.Anything, mock.MatchedBy(func(rawToken string) bool {
					return rawToken == rt
				})).Return(nil, tokenerrors.ErrNotFound)
			},
			wantErr: tokenerrors.ErrNotFound,
			want:    0,
		},

		{
			name:     "Ошибка БД",
			rawToken: "rnd-token",
			setup: func(tokenService *mocks.TokenServiceMock, rt string) {
				tokenService.EXPECT().Get(mock.Anything, mock.MatchedBy(func(rawToken string) bool {
					return rawToken == rt
				})).Return(nil, errDB)
			},
			wantErr: errDB,
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenService := mocks.NewTokenServiceMock(t)
			userSvc := mocks.NewUserServiceMock(t)
			transactor := mocks.NewTransactorMock(t)
			logger := log.New(io.Discard, "", 0)

			tt.setup(tokenService, tt.rawToken)

			h := auth.NewHandler(userSvc, tokenService, logger, transactor)

			got, gotErr := h.Authorize(context.Background(), tt.rawToken)

			if tt.wantErr != nil {
				require.ErrorIs(t, gotErr, tt.wantErr)
			} else {
				require.NoError(t, gotErr)
			}

			require.Equal(t, got, tt.want)
		})
	}
}
