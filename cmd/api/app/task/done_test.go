package api_test

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

	"todoshnik/cmd/api/app/task/mocks"
	"todoshnik/internal/domains/task"
	"todoshnik/internal/domains/user"
	"todoshnik/internal/infrastructure/identity"

	api "todoshnik/cmd/api/app/task"
	taskerrors "todoshnik/internal/domains/task/errors"
	usererrors "todoshnik/internal/domains/user/errors"
	authcontext "todoshnik/internal/infrastructure/context_manager/auth"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Done(t *testing.T) {
	testScope := identity.AccessScope{
		UserID: 10,
	}
	testTaskID := 1

	tests := []struct {
		name     string
		taskID   int
		withUser bool
		setup    func(
			service *mocks.TaskServiceMock,
			userSvc *mocks.UserGetterMock,
		)
		wantStatus int
		wantTask   *task.Task
		taskURI    string
	}{
		{
			name:     "Успешно изменен статус",
			withUser: true,
			taskID:   testTaskID,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					MarkDone(mock.Anything, testTaskID, testScope).
					Return(&task.Task{
						ID:     1,
						Title:  "Тестовый таск",
						Done:   true,
						UserID: testScope.UserID,
					}, nil)
			},
			wantTask: &task.Task{
				ID:     1,
				Title:  "Тестовый таск",
				Done:   true,
				UserID: testScope.UserID,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "Задача не найдена",
			withUser: true,
			taskID:   testTaskID,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					MarkDone(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, taskerrors.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "ошибка сервиса",
			taskID:   testTaskID,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					MarkDone(mock.Anything, 1, testScope).
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "нет пользователя в контексте",
			withUser:   false,
			taskID:     testTaskID,
			setup:      func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "пользователь не найден",
			withUser: true,
			taskID:   testTaskID,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, mock.AnythingOfType("int")).
					Return(nil, usererrors.ErrNotFound)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:    "невалидный ID задачи",
			taskURI: "abc",
			setup: func(
				service *mocks.TaskServiceMock,
				userSvc *mocks.UserGetterMock,
			) {
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewTaskServiceMock(t)
			userSvc := mocks.NewUserGetterMock(t)
			logger := log.New(io.Discard, "", 0)

			tt.setup(service, userSvc)

			handler := api.NewHandler(service, userSvc, logger)

			router := chi.NewRouter()
			router.Post(
				"/tasks/{id}",
				handler.Done,
			)
			reqPath := "/tasks/" + tt.taskURI
			if tt.taskURI == "" {
				reqPath = "/tasks/" + strconv.Itoa(tt.taskID)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				reqPath,
				strings.NewReader(""),
			)

			if tt.withUser {
				req = req.WithContext(
					authcontext.SetUserID(req.Context(), testScope.UserID),
				)
			}

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if tt.wantTask != nil {
				var got *task.Task

				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
				require.Equal(t, tt.wantTask, got)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
