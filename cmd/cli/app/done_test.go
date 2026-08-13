package cli

import (
	"context"
	"errors"
	"testing"

	"todoshnik/cmd/cli/app/mocks"
	"todoshnik/internal/domains/task"
	"todoshnik/internal/infrastructure/identity"

	"github.com/stretchr/testify/mock"
)

func TestHandler_done(t *testing.T) {
	errDB := errors.New("ошибка БД")

	testTask := &task.Task{
		ID:     15,
		Title:  "Тестовая задача",
		Done:   true,
		UserID: 10,
	}

	tests := []struct {
		name  string
		args  []string
		setup func(taskService *mocks.TaskServiceMock)
	}{
		{
			name: "не указан ID",
			args: []string{"cli", "done"},
			setup: func(taskService *mocks.TaskServiceMock) {
			},
		},
		{
			name: "некорректный ID",
			args: []string{"cli", "done", "abc"},
			setup: func(taskService *mocks.TaskServiceMock) {
			},
		},
		{
			name: "ошибка пометки задачи",
			args: []string{"cli", "done", "15"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					MarkDone(
						mock.Anything,
						15,
						identity.AccessScope{IsAdmin: true},
					).
					Return(nil, errDB)
			},
		},
		{
			name: "успех",
			args: []string{"cli", "done", "15"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					MarkDone(
						mock.Anything,
						15,
						identity.AccessScope{IsAdmin: true},
					).
					Return(testTask, nil)
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

			handler.done(context.Background(), tt.args)
		})
	}
}
