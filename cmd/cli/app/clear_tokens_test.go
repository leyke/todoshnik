package cli

import (
	"context"
	"errors"
	"testing"

	"todoshnik/cmd/cli/app/mocks"

	"github.com/stretchr/testify/mock"
)

func TestHandler_clearTokens(t *testing.T) {
	errDB := errors.New("ошибка БД")

	tests := []struct {
		name  string
		setup func(tokenService *mocks.TokenServiceMock)
	}{
		{
			name: "успех",
			setup: func(tokenService *mocks.TokenServiceMock) {
				tokenService.EXPECT().
					ClearExpiredTokens(mock.Anything).
					Return(3, nil)
			},
		},
		{
			name: "ошибка удаления токенов",
			setup: func(tokenService *mocks.TokenServiceMock) {
				tokenService.EXPECT().
					ClearExpiredTokens(mock.Anything).
					Return(0, errDB)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenService := mocks.NewTokenServiceMock(t)

			tt.setup(tokenService)

			handler := &Handler{
				tokenService: tokenService,
			}

			handler.clearTokens(context.Background())
		})
	}
}
