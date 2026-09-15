package cli

import (
	"context"
	"errors"
	"testing"

	"todoshnik/cmd/cli/app/mocks"
	"todoshnik/internal/infrastructure/identity"

	"github.com/stretchr/testify/mock"
)

func TestHandler_delete(t *testing.T) {
	errDB := errors.New("ошибка БД")

	tests := []struct {
		name  string
		args  []string
		setup func(taskService *mocks.TaskServiceMock)
	}{
		{
			name: "не указан ID",
			args: []string{"cli", "delete"},
			setup: func(taskService *mocks.TaskServiceMock) {
			},
		},
		{
			name: "некорректный ID",
			args: []string{"cli", "delete", "abc"},
			setup: func(taskService *mocks.TaskServiceMock) {
			},
		},
		{
			name: "ошибка удаления задачи",
			args: []string{"cli", "delete", "15"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					Delete(
						mock.Anything,
						15,
						identity.AccessScope{IsAdmin: true},
					).
					Return(errDB)
			},
		},
		{
			name: "успех",
			args: []string{"cli", "delete", "15"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					Delete(
						mock.Anything,
						15,
						identity.AccessScope{IsAdmin: true},
					).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskService := mocks.NewTaskServiceMock(t)

			tt.setup(taskService)

			handler := &Handler{
				taskService: taskService,
			}

			handler.delete(context.Background(), tt.args)
		})
	}
}
