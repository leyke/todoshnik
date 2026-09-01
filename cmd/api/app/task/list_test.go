package api_test

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todoshnik/cmd/api/app/task/mocks"
	"todoshnik/internal/domains/task"
	"todoshnik/internal/domains/user"
	"todoshnik/internal/infrastructure/identity"

	api "todoshnik/cmd/api/app/task"
	usererrors "todoshnik/internal/domains/user/errors"
	authcontext "todoshnik/internal/infrastructure/context_manager/auth"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_List(t *testing.T) {
	testScope := identity.AccessScope{
		UserID: 10,
	}

	testTask := &task.Task{
		ID:     1,
		UserID: testScope.UserID,
		Done:   true,
		Title:  "Test",
	}

	tests := []struct {
		name     string
		withUser bool
		setup    func(
			service *mocks.TaskServiceMock,
			userSvc *mocks.UserGetterMock,
			count int,
		)
		wantCount  int
		wantStatus int
	}{
		{
			name:     "нет пользователя в контексте",
			withUser: false,
			setup: func(
				service *mocks.TaskServiceMock,
				userSvc *mocks.UserGetterMock,
				count int,
			) {
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "нет пользователя в бд",
			withUser: true,
			setup: func(
				service *mocks.TaskServiceMock,
				userSvc *mocks.UserGetterMock,
				count int,
			) {
				userSvc.EXPECT().GetById(mock.Anything, testScope.UserID).Return(nil, usererrors.ErrNotFound)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "нет задач",
			withUser: true,
			setup: func(
				service *mocks.TaskServiceMock,
				userSvc *mocks.UserGetterMock,
				count int,
			) {
				userSvc.EXPECT().GetById(mock.Anything, testScope.UserID).Return(&user.User{ID: 1}, nil)

				service.EXPECT().List(mock.Anything, task.TaskFilter{Scope: testScope}).Return(make([]*task.Task, 0), nil)
			},
			wantCount:  0,
			wantStatus: http.StatusOK,
		},
		{
			name:     "1 задача",
			withUser: true,
			setup: func(
				service *mocks.TaskServiceMock,
				userSvc *mocks.UserGetterMock,
				count int,
			) {
				userSvc.EXPECT().GetById(mock.Anything, testScope.UserID).Return(&user.User{ID: 1}, nil)

				var items []*task.Task
				for range count {
					testTask.ID += 1
					testTask.Done = testTask.ID/2 == 0
					items = append(items, testTask)
				}

				service.EXPECT().List(mock.Anything, task.TaskFilter{Scope: testScope}).Return(items, nil)
			},
			wantCount:  1,
			wantStatus: http.StatusOK,
		},
		{
			name:     "5 задач",
			withUser: true,
			setup: func(
				service *mocks.TaskServiceMock,
				userSvc *mocks.UserGetterMock,
				count int,
			) {
				userSvc.EXPECT().GetById(mock.Anything, testScope.UserID).Return(&user.User{ID: 1}, nil)

				var items []*task.Task
				for range count {
					testTask.ID += 1
					testTask.Done = testTask.ID/2 == 0
					items = append(items, testTask)
				}

				service.EXPECT().List(mock.Anything, task.TaskFilter{Scope: testScope}).Return(items, nil)
			},
			wantCount:  5,
			wantStatus: http.StatusOK,
		},
		{
			name:     "ошибка бд",
			withUser: true,
			setup: func(
				service *mocks.TaskServiceMock,
				userSvc *mocks.UserGetterMock,
				count int,
			) {
				userSvc.EXPECT().GetById(mock.Anything, testScope.UserID).Return(&user.User{ID: 1}, nil)

				service.EXPECT().List(mock.Anything, task.TaskFilter{Scope: testScope}).Return(nil, errors.New("Ошибка БД"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			service := mocks.NewTaskServiceMock(t)
			userSvc := mocks.NewUserGetterMock(t)
			logger := log.New(io.Discard, "", 0)

			tt.setup(service, userSvc, tt.wantCount)

			handler := api.NewHandler(service, userSvc, logger)

			req := httptest.NewRequest(
				http.MethodGet,
				"/tasks",
				strings.NewReader(""),
			)
			if tt.withUser {
				req = req.WithContext(
					authcontext.SetUserID(req.Context(), testScope.UserID),
				)
			}

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.List(rec, req)

			if tt.wantCount != 0 {
				var got []*task.Task
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

				require.Equal(t, tt.wantCount, len(got))
			}

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
