package api_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"todoshnik/cmd/api/app/task/mocks"
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

func TestHandler_Delete(t *testing.T) {
	type testCase struct {
		name     string
		taskID   int
		withUser bool
		setup    func(
			service *mocks.TaskServiceMock,
			userSvc *mocks.UserGetterMock,
		)
		wantStatus int
		taskURI    string
	}

	testScope := identity.AccessScope{
		UserID: 10,
	}
	testTaskID := 1

	tests := []testCase{
		{
			name:     "Успешно удален",
			withUser: true,
			taskID:   testTaskID,
			setup: func(service *mocks.TaskServiceMock, userSvc *mocks.UserGetterMock) {
				userSvc.EXPECT().
					GetById(mock.Anything, testScope.UserID).
					Return(&user.User{ID: testScope.UserID}, nil)

				service.EXPECT().
					Delete(mock.Anything, testTaskID, testScope).
					Return(nil)
			},
			wantStatus: http.StatusNoContent,
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
					Delete(mock.Anything, mock.Anything, mock.Anything).
					Return(taskerrors.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
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
			router.Delete(
				"/tasks/{id}",
				handler.Delete,
			)
			reqPath := "/tasks/" + tt.taskURI
			if tt.taskURI == "" {
				reqPath = "/tasks/" + strconv.Itoa(tt.taskID)
			}

			req := httptest.NewRequest(
				http.MethodDelete,
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

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
