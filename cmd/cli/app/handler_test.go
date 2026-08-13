package cli

import (
	"testing"

	"todoshnik/cmd/cli/app/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Run(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		setup func(
			taskService *mocks.TaskServiceMock,
			tokenService *mocks.TokenServiceMock,
		)
	}{
		{
			name: "нет команды",
			args: []string{"app"},
			setup: func(
				taskService *mocks.TaskServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
			},
		},
		{
			name: "add",
			args: []string{"app", "add"},
			setup: func(
				taskService *mocks.TaskServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
			},
		},
		{
			name: "list",
			args: []string{"app", "list"},
			setup: func(
				taskService *mocks.TaskServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
				taskService.EXPECT().
					List(mock.Anything, mock.Anything).
					Return(nil, nil)
			},
		},
		{
			name: "done",
			args: []string{"app", "done"},
			setup: func(
				taskService *mocks.TaskServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
			},
		},
		{
			name: "delete",
			args: []string{"app", "delete"},
			setup: func(
				taskService *mocks.TaskServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
			},
		},
		{
			name: "clear-tokens",
			args: []string{"app", "clear-tokens"},
			setup: func(
				taskService *mocks.TaskServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
				tokenService.EXPECT().
					ClearExpiredTokens(mock.Anything).
					Return(0, nil)
			},
		},
		{
			name: "неизвестная команда",
			args: []string{"app", "unknown"},
			setup: func(
				taskService *mocks.TaskServiceMock,
				tokenService *mocks.TokenServiceMock,
			) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskService := mocks.NewTaskServiceMock(t)
			tokenService := mocks.NewTokenServiceMock(t)

			tt.setup(taskService, tokenService)

			handler := NewHandler(taskService, tokenService)

			require.NotPanics(t, func() {
				handler.Run(tt.args)
			})
		})
	}
}
