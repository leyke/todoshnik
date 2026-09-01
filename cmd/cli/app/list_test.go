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

func TestHandler_list(t *testing.T) {
	errDB := errors.New("ошибка БД")

	testTasks := []*task.Task{
		{
			ID:     1,
			Title:  "Первая задача",
			Done:   false,
			UserID: 10,
		},
		{
			ID:     2,
			Title:  "Вторая задача",
			Done:   true,
			UserID: 10,
		},
	}

	tests := []struct {
		name  string
		args  []string
		setup func(taskService *mocks.TaskServiceMock)
	}{
		{
			name: "получение всех задач",
			args: []string{"cli", "list"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					List(
						mock.Anything,
						task.TaskFilter{
							Status: "",
							Scope:  identity.AccessScope{IsAdmin: true},
						},
					).
					Return(testTasks, nil)
			},
		},
		{
			name: "фильтр по статусу",
			args: []string{"cli", "list", "-status", "completed"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					List(
						mock.Anything,
						task.TaskFilter{
							Status: task.StatusCompleted,
							Scope:  identity.AccessScope{IsAdmin: true},
						},
					).
					Return(testTasks, nil)
			},
		},
		{
			name: "пустой список",
			args: []string{"cli", "list"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					List(
						mock.Anything,
						task.TaskFilter{
							Status: "",
							Scope:  identity.AccessScope{IsAdmin: true},
						},
					).
					Return([]*task.Task{}, nil)
			},
		},
		{
			name: "ошибка получения списка",
			args: []string{"cli", "list"},
			setup: func(taskService *mocks.TaskServiceMock) {
				taskService.EXPECT().
					List(
						mock.Anything,
						task.TaskFilter{
							Status: "",
							Scope:  identity.AccessScope{IsAdmin: true},
						},
					).
					Return(nil, errDB)
			},
		},
		{
			name: "ошибка получения параметров",
			args: []string{"app", "list", "-status"},
			setup: func(taskService *mocks.TaskServiceMock) {
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

			handler.list(context.Background(), tt.args)
		})
	}
}
