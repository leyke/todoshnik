package cli

import (
	"context"
	"errors"
	"testing"
	
	"todoshnik/cmd/cli/app/mocks"
	"todoshnik/internal/domains/task"

	"github.com/stretchr/testify/mock"
)

func TestHandler_add(t *testing.T) {
	errDB := errors.New("ошибка БД")

	testTask := &task.Task{
		ID:    15,
		Title: "Новая задача",
	}

	tests := []struct {
		name  string
		args  []string
		setup func(taskService *mocks.TaskServiceMock)
	}{
		{
			name: "не указано название задачи",
			args: []string{"cli", "add"},
			setup: func(taskService *mocks.TaskServiceMock) {
			},
		},
		{
			name: "ошибка ID",
			args: []string{"cli", "add", "Новая задача", "abc"},
			setup: func(taskService *mocks.TaskServiceMock) {
			},
		},
		{
			name: "ошибка создания задачи",
			args: []string{"cli", "add", "Новая задача", "10"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					Add(mock.Anything, "Новая задача", 10).
					Return(nil, errDB)
			},
		},
		{
			name: "успех",
			args: []string{"cli", "add", "Новая задача", "10"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					Add(mock.Anything, "Новая задача", 10).
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

			handler.add(context.Background(), tt.args)

			taskService.AssertExpectations(t)
		})
	}
}
