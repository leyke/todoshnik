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
	authcontext "todoshnik/internal/infrastructure/context_manager/auth"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_View(t *testing.T) {
	testScope := identity.AccessScope{
		UserID: 10,
	}
	testTitle := "покрыть тестами HTTP Update"
	testTaskId := 1

	expect := &task.Task{
		ID:     testTaskId,
		Title:  testTitle,
		UserID: testScope.UserID,
		Done:   false,
	}

	tests := []struct {
		name     string
		withUser bool
		setup    func(
			service *mocks.TaskServiceMock,
			userSvc *mocks.UserGetterMock,
		)
		wantStatus int
		wantTask   *task.Task
		wantTaskID int
		taskURI    string
	}{
		{
			name:     "Задача не найдена",
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					Get(mock.Anything, mock.Anything, testScope).
					Return(nil, taskerrors.ErrNotFound)
			},
			wantTaskID: expect.ID,
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "успешное изменение",
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					Get(mock.Anything, expect.ID, testScope).
					Return(expect, nil)
			},
			wantTaskID: expect.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "нет пользователя в контексте",
			withUser:   false,
			setup:      func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {},
			wantTaskID: expect.ID,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "пользователь не найден",
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(nil, errors.New("not found"))
			},
			wantTaskID: expect.ID,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "ошибка сервиса",
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					Get(mock.Anything, expect.ID, testScope).
					Return(nil, errors.New("db error"))
			},
			wantTaskID: expect.ID,
			wantStatus: http.StatusInternalServerError,
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
			router.Get(
				"/tasks/{id}",
				handler.View,
			)

			reqPath := "/tasks/" + tt.taskURI
			if tt.taskURI == "" {
				reqPath = "/tasks/" + strconv.Itoa(tt.wantTaskID)
			}

			req := httptest.NewRequest(
				http.MethodGet,
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
