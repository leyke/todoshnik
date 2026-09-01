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
	"todoshnik/internal/infrastructure/validation"

	api "todoshnik/cmd/api/app/task"
	authcontext "todoshnik/internal/infrastructure/context_manager/auth"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Create(t *testing.T) {
	type testCase struct {
		name     string
		body     string
		withUser bool
		setup    func(
			service *mocks.TaskServiceMock,
			userSvc *mocks.UserGetterMock,
		)
		wantStatus int
		wantTask   *task.Task
	}

	testUserID := 10
	testTitle := "покрыть тестами HTTP Create"

	expect := &task.Task{
		ID:     1,
		Title:  testTitle,
		UserID: testUserID,
	}

	tests := []testCase{
		{
			name:     "успешное создание",
			body:     `{"title":"` + testTitle + `"}`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testUserID).
					Return(&user.User{ID: testUserID}, nil)

				service.EXPECT().
					Add(mock.Anything, testTitle, testUserID).
					Return(&task.Task{
						ID:     1,
						Title:  testTitle,
						UserID: testUserID,
						
					}, nil)
			},
			wantTask:   expect,
			wantStatus: http.StatusOK,
		},
		{
			name:       "нет пользователя в контексте",
			body:       `{"title":"` + testTitle + `"}`,
			withUser:   false,
			setup:      func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "пользователь не найден",
			body:     `{"title":"` + testTitle + `"}`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testUserID).
					Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "ошибка сервиса",
			body:     `{"title":"` + testTitle + `"}`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testUserID).
					Return(&user.User{ID: testUserID}, nil)

				service.EXPECT().
					Add(mock.Anything, testTitle, testUserID).
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:     "ошибка валидации",
			body:     `{"title":""}`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testUserID).
					Return(&user.User{ID: testUserID}, nil)

				service.EXPECT().
					Add(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, validation.NewValidationErrorFromValidator(nil))
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:     "невалидный json",
			body:     `{"title":"`,
			withUser: true,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "нет пользователя в контексте",
			body: `{"title":"` + testTitle + `"}`,
			setup: func(
				service *mocks.TaskServiceMock,
				userSvc *mocks.UserGetterMock,
			) {
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewTaskServiceMock(t)
			userSvc := mocks.NewUserGetterMock(t)
			logger := log.New(io.Discard, "", 0)

			tt.setup(service, userSvc)

			handler := api.NewHandler(service, userSvc, logger)

			req := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				strings.NewReader(tt.body),
			)

			if tt.withUser {
				req = req.WithContext(
					authcontext.SetUserID(req.Context(), testUserID),
				)
			}

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			if tt.wantTask != nil {
				var got *task.Task
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

				require.Equal(t, tt.wantTask, got)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
