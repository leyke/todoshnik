package task_test

import (
	"context"
	"errors"
	"testing"

	"todoshnik/internal/domains/task"
	"todoshnik/internal/domains/task/mocks"
	"todoshnik/internal/infrastructure/identity"

	taskerrors "todoshnik/internal/domains/task/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Add(t *testing.T) {
	ctx := context.Background()
	errDB := errors.New("ошибка БД")

	tests := []struct {
		name    string
		title   string
		userID  int
		setup   func(repo *mocks.RepositoryMock)
		wantErr error
	}{
		{
			name:   "успешное создание задачи",
			title:  "Прикрутить тесты к проекту",
			userID: 10,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					Create(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, t *task.Task) (*task.Task, error) {
						t.ID = 1
						return t, nil
					})
			},
		},
		{
			name:   "ошибка репозитория",
			title:  "Прикрутить тесты к проекту",
			userID: 10,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					Create(
						mock.Anything,
						mock.Anything,
					).
					Return(
						nil,
						errDB,
					)
			},
			wantErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepositoryMock(t)

			tt.setup(repo)

			service := task.NewService(repo)

			result, err := service.Add(
				ctx,
				tt.title,
				tt.userID,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.title, result.Title)
		})
	}
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()

	testScope := identity.AccessScope{
		UserID: 10,
	}

	newTestTitle := "Прикрутить еще тесты"
	newTestDone := true

	existingTestTask := &task.Task{
		ID:     15,
		Title:  "Старая задача",
		Done:   false,
		UserID: testScope.UserID,
	}

	errDB := errors.New("ошибка БД")

	tests := []struct {
		name    string
		title   string
		done    bool
		scope   identity.AccessScope
		setup   func(repo *mocks.RepositoryMock)
		wantErr error
	}{
		{
			name:  "успешное обновление задачи",
			title: newTestTitle,
			done:  newTestDone,
			scope: testScope,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					GetByID(
						mock.Anything,
						existingTestTask.ID,
						testScope,
					).
					Return(existingTestTask, nil)

				repo.EXPECT().
					Update(
						mock.Anything,
						mock.MatchedBy(func(t *task.Task) bool {
							return t.ID == existingTestTask.ID &&
								t.UserID == existingTestTask.UserID &&
								t.Title == newTestTitle &&
								t.Done == newTestDone
						}),
					).
					Return(nil)
			},
		},
		{
			name:  "ошибка репозитория при обновлении",
			title: newTestTitle,
			done:  newTestDone,
			scope: testScope,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					GetByID(
						mock.Anything,
						existingTestTask.ID,
						testScope,
					).
					Return(existingTestTask, nil)

				repo.EXPECT().
					Update(
						mock.Anything,
						mock.MatchedBy(func(t *task.Task) bool {
							return t.ID == existingTestTask.ID &&
								t.UserID == existingTestTask.UserID &&
								t.Title == newTestTitle &&
								t.Done == newTestDone
						}),
					).
					Return(errDB)
			},
			wantErr: errDB,
		},
		{
			name:  "ошибка при получении задачи",
			title: newTestTitle,
			done:  newTestDone,
			scope: testScope,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					GetByID(
						mock.Anything,
						existingTestTask.ID,
						testScope,
					).
					Return(nil, errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepositoryMock(t)

			tt.setup(repo)

			service := task.NewService(repo)

			result, err := service.Update(
				ctx,
				existingTestTask.ID,
				tt.title,
				tt.done,
				tt.scope,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			expected := &task.Task{
				ID:     existingTestTask.ID,
				UserID: existingTestTask.UserID,
				Title:  tt.title,
				Done:   tt.done,
			}

			require.Equal(t, expected, result)
		})
	}
}

func TestService_Get(t *testing.T) {
	testScope := identity.AccessScope{
		UserID: 10,
	}

	expected := &task.Task{
		ID:     15,
		Title:  "Старая задача",
		Done:   false,
		UserID: testScope.UserID,
	}

	errDB := errors.New("ошибка БД")

	tests := []struct {
		name    string
		setup   func(repo *mocks.RepositoryMock)
		taskId  int
		scope   identity.AccessScope
		want    *task.Task
		wantErr error
	}{
		{
			name:   "успех при получении задачи",
			scope:  testScope,
			taskId: expected.ID,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					GetByID(
						mock.Anything,
						expected.ID,
						testScope,
					).
					Return(expected, nil)
			},
			want: expected,
		},
		{
			name:   "ошибка бд при получении задачи",
			scope:  testScope,
			taskId: expected.ID,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					GetByID(
						mock.Anything,
						expected.ID,
						testScope,
					).
					Return(nil, errDB)
			},
			want:    nil,
			wantErr: errDB,
		},
		{
			name:   "задача не найдена",
			scope:  testScope,
			taskId: expected.ID,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					GetByID(
						mock.Anything,
						expected.ID,
						testScope,
					).
					Return(nil, taskerrors.ErrNotFound)
			},
			want:    nil,
			wantErr: taskerrors.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepositoryMock(t)

			tt.setup(repo)

			service := task.NewService(repo)

			got, gotErr := service.Get(context.Background(), tt.taskId, tt.scope)

			if tt.wantErr != nil {
				require.ErrorIs(t, gotErr, tt.wantErr)
				return
			}

			require.NoError(t, gotErr)
			require.Equal(t, expected, got)
		})
	}
}

func TestService_MarkDone(t *testing.T) {
	testScope := identity.AccessScope{
		UserID: 10,
	}

	existingTestTask := &task.Task{
		ID:     15,
		Title:  "Невыполненная задача",
		Done:   false,
		UserID: testScope.UserID,
	}

	existingDonedTask := &task.Task{
		ID:     16,
		Title:  "Выполненная задача",
		Done:   true,
		UserID: testScope.UserID,
	}

	testUndone := false
	testDone := true
	errDB := errors.New("ошибка БД")

	tests := []struct {
		name     string
		setup    func(repo *mocks.RepositoryMock)
		taskId   int
		scope    identity.AccessScope
		wantDone bool
		wantErr  error
	}{
		{
			name:     "успех отметки - выполненна",
			taskId:   existingTestTask.ID,
			wantDone: testDone,
			scope:    testScope,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					GetByID(
						mock.Anything,
						existingTestTask.ID,
						testScope,
					).
					Return(existingTestTask, nil)

				repo.EXPECT().
					Update(
						mock.Anything,
						mock.MatchedBy(func(t *task.Task) bool {
							return t.ID == existingTestTask.ID &&
								t.UserID == existingTestTask.UserID &&
								t.Title == existingTestTask.Title &&
								t.Done == testDone
						}),
					).
					Return(nil)
			},
		},
		{
			name:     "успех отметки - невыполнена",
			taskId:   existingDonedTask.ID,
			wantDone: testUndone,
			scope:    testScope,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					GetByID(
						mock.Anything,
						existingDonedTask.ID,
						testScope,
					).
					Return(existingDonedTask, nil)

				repo.EXPECT().
					Update(
						mock.Anything,
						mock.MatchedBy(func(t *task.Task) bool {
							return t.ID == existingDonedTask.ID &&
								t.UserID == existingDonedTask.UserID &&
								t.Title == existingDonedTask.Title &&
								t.Done == testUndone
						}),
					).
					Return(nil)
			},
		},

		{
			name:     "ошибка репозитория",
			taskId:   existingTestTask.ID,
			wantDone: testUndone,
			scope:    testScope,
			wantErr:  errDB,
			setup: func(repo *mocks.RepositoryMock) {
				repo.EXPECT().
					GetByID(
						mock.Anything,
						existingTestTask.ID,
						testScope,
					).
					Return(nil, errDB)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepositoryMock(t)

			tt.setup(repo)

			service := task.NewService(repo)

			got, gotErr := service.MarkDone(context.Background(), tt.taskId, tt.scope)

			if tt.wantErr != nil {
				require.ErrorIs(t, gotErr, tt.wantErr)
				return
			}

			require.NoError(t, gotErr)
			require.Equal(t, tt.wantDone, got.Done)
		})
	}
}
