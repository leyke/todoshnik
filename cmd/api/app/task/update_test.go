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
	"todoshnik/internal/infrastructure/validation"

	api "todoshnik/cmd/api/app/task"
	taskerrors "todoshnik/internal/domains/task/errors"
	authcontext "todoshnik/internal/infrastructure/context_manager/auth"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Update(t *testing.T) {
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
		body     string
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
			body:     `{"title":"` + expect.Title + `"}`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					Update(mock.Anything, mock.Anything, mock.Anything, mock.Anything, testScope).
					Return(nil, taskerrors.ErrNotFound)
			},
			wantTaskID: expect.ID,
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "успешное изменение",
			body:     `{"title":"` + expect.Title + `"}`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					Update(mock.Anything, expect.ID, expect.Title, expect.Done, testScope).
					Return(expect, nil)
			},
			wantTaskID: expect.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "нет пользователя в контексте",
			body:       `{"title":"` + testTitle + `"}`,
			withUser:   false,
			setup:      func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {},
			wantTaskID: expect.ID,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "пользователь не найден",
			body:     `{"title":"` + testTitle + `"}`,
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
			body:     `{"title":"` + testTitle + `"}`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					Update(mock.Anything, expect.ID, expect.Title, expect.Done, testScope).
					Return(nil, errors.New("db error"))
			},
			wantTaskID: expect.ID,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:     "ошибка валидации",
			body:     `{"title":""}`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					Update(mock.Anything, expect.ID, "", expect.Done, testScope).
					Return(nil, validation.NewValidationErrorFromValidator(nil))
			},
			wantTaskID: expect.ID,
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:     "невалидный json",
			body:     `{"title":"`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
			},
			wantTaskID: expect.ID,
			wantStatus: http.StatusBadRequest,
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
				handler.Update,
			)
			reqPath := "/tasks/" + tt.taskURI
			if tt.taskURI == "" {
				reqPath = "/tasks/" + strconv.Itoa(tt.wantTaskID)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				reqPath,
				strings.NewReader(tt.body),
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
