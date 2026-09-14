package user_test

import (
	"context"
	"errors"
	"testing"

	"todoshnik/internal/domains/user"
	"todoshnik/internal/domains/user/mocks"
	"todoshnik/internal/infrastructure/utils/test"
	"todoshnik/internal/infrastructure/validation"

	usererrors "todoshnik/internal/domains/user/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Add(t *testing.T) {
	ctx := context.Background()

	errDB := errors.New("ошибка БД")
	errHash := errors.New("ошибка хеширования")

	tests := []struct {
		name     string
		userName string
		login    string
		password string
		setup    func(
			repo *mocks.RepositoryMock,
			hasher *test.FakePwdHasher,
		)
		wantErr error
	}{
		{
			name:     "успешное создание",
			userName: "Пользователь",
			login:    "login",
			password: "password",
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByLogin(mock.Anything, "login").
					Return(nil, nil)

				repo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(&user.User{
						ID:           1,
						Name:         "Пользователь",
						Login:        "login",
						PasswordHash: "hashed_password",
					}, nil)
			},
		},
		{
			name:     "пользователь уже существует",
			userName: "Пользователь",
			login:    "login",
			password: "password",
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				existingUser := &user.User{
					ID:   10,
					Name: "Другой пользователь",
				}

				repo.EXPECT().
					GetByLogin(mock.Anything, "login").
					Return(existingUser, nil)
			},
			wantErr: usererrors.ErrConflict,
		},
		{
			name:     "ошибка хеширования пароля",
			userName: "Пользователь",
			login:    "login",
			password: "password",
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByLogin(mock.Anything, "login").
					Return(nil, nil)

				hasher.HashErr = errHash
			},
			wantErr: errHash,
		},
		{
			name:     "ошибка репозитория (создания)",
			userName: "Пользователь",
			login:    "login",
			password: "password",
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByLogin(mock.Anything, "login").
					Return(nil, nil)

				repo.EXPECT().
					Create(
						mock.Anything,
						mock.Anything,
					).
					Return(nil, errDB)
			},
			wantErr: errDB,
		},

		{
			name:     "ошибка репозитория (поиска существуюшего)",
			userName: "Пользователь",
			login:    "login",
			password: "password",
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByLogin(mock.Anything, "login").
					Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name:     "невалидное имя пользователя",
			userName: "a",
			login:    "login",
			password: "password",
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByLogin(mock.Anything, "login").
					Return(nil, nil)
			},
			wantErr: validation.ErrNotValidate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepositoryMock(t)
			hasher := test.NewFakePwdHasher()

			tt.setup(repo, hasher)

			service := user.NewService(repo, hasher)

			result, err := service.Add(
				ctx,
				tt.userName,
				tt.login,
				tt.password,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			require.Equal(t, tt.userName, result.Name)
			require.Equal(t, tt.login, result.Login)
			require.Equal(t, "hashed_password", result.PasswordHash)

			require.Equal(t, 1, hasher.HashCalls)
			require.Equal(t, tt.password, hasher.LastHashInput)
		})
	}
}

func TestService_AddFromTg(t *testing.T) {
	ctx := context.Background()

	errDB := errors.New("ошибка БД")

	tests := []struct {
		name       string
		userName   string
		telegramID int64
		setup      func(
			repo *mocks.RepositoryMock,
			hasher *test.FakePwdHasher,
		)
		wantErr error
	}{
		{
			name:       "успешное создание",
			userName:   "Пользователь",
			telegramID: int64(12345),
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByTgId(mock.Anything, int64(12345)).
					Return(nil, nil)

				repo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(&user.User{
						ID:         1,
						Name:       "Пользователь",
						TelegramID: 12345,
					}, nil)
			},
		},
		{
			name:       "пользователь уже существует",
			userName:   "Пользователь",
			telegramID: int64(12345),
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				existingUser := &user.User{
					ID:         1,
					Name:       "Пользователь",
					TelegramID: 12345,
				}

				repo.EXPECT().
					GetByTgId(mock.Anything, int64(12345)).
					Return(existingUser, nil)
			},
		},
		{
			name:       "ошибка репозитория (создания)",
			userName:   "Пользователь",
			telegramID: 12345,
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByTgId(mock.Anything, int64(12345)).
					Return(nil, nil)

				repo.EXPECT().
					Create(
						mock.Anything,
						mock.Anything,
					).
					Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name:       "ошибка репозитория (поиска существуюшего)",
			userName:   "Пользователь",
			telegramID: int64(12345),
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByTgId(mock.Anything, int64(12345)).
					Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name:       "невалидное имя пользователя",
			userName:   "a",
			telegramID: int64(12345),
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByTgId(mock.Anything, int64(12345)).
					Return(nil, nil)
			},
			wantErr: validation.ErrNotValidate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepositoryMock(t)
			hasher := test.NewFakePwdHasher()

			tt.setup(repo, hasher)

			service := user.NewService(repo, hasher)

			result, err := service.AddFromTg(
				ctx,
				tt.userName,
				tt.telegramID,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			require.Equal(t, tt.telegramID, result.TelegramID)

			require.Equal(t, 0, hasher.HashCalls)
		})
	}
}

func TestService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("successful", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t)
		hasher := test.NewFakePwdHasher()

		expectedUsers := []*user.User{
			{
				ID:    1,
				Name:  "Artem",
				Login: "123456789",
			},
			{
				ID:    2,
				Name:  "Ivan",
				Login: "987654321",
			},
		}

		repo.EXPECT().
			List(mock.Anything).
			Return(expectedUsers, nil)

		service := user.NewService(repo, hasher)

		got, err := service.List(ctx)

		require.NoError(t, err)
		require.Equal(t, expectedUsers, got)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t)
		hasher := test.NewFakePwdHasher()

		expectedErr := errors.New("ошибка БД")

		repo.EXPECT().
			List(mock.Anything).
			Return(nil, expectedErr)

		service := user.NewService(repo, hasher)

		got, err := service.List(ctx)

		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, got)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()

	errDB := errors.New("ошибка БД")

	existingUser := &user.User{
		ID:           1,
		Name:         "Пользователь",
		Login:        "login",
		PasswordHash: "hashed_password",
	}
	tests := []struct {
		name     string
		userName string
		userID   int
		setup    func(
			repo *mocks.RepositoryMock,
			hasher *test.FakePwdHasher,
		)
		wantErr error
	}{
		{
			name:     "успешное обновление",
			userName: "Пользователь UPD",
			userID:   1,
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByID(mock.Anything, 1).
					Return(existingUser, nil)

				repo.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(nil)
			},
		},
		{
			name:     "ошибка репозитория (обновлеия)",
			userName: "Пользователь",
			userID:   1,
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByID(mock.Anything, 1).
					Return(existingUser, nil)

				repo.EXPECT().
					Update(
						mock.Anything,
						mock.Anything,
					).
					Return(errDB)
			},
			wantErr: errDB,
		},
		{
			name:     "ошибка репозитория (поиска существуюшего)",
			userName: "Пользователь",
			userID:   1,
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByID(mock.Anything, 1).
					Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name:     "невалидное имя пользователя",
			userName: "a",
			userID:   1,
			setup: func(repo *mocks.RepositoryMock, hasher *test.FakePwdHasher) {
				repo.EXPECT().
					GetByID(mock.Anything, 1).
					Return(existingUser, nil)
			},
			wantErr: validation.ErrNotValidate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepositoryMock(t)
			hasher := test.NewFakePwdHasher()

			tt.setup(repo, hasher)

			service := user.NewService(repo, hasher)

			result, err := service.Update(
				ctx,
				tt.userID,
				tt.userName,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			require.Equal(t, tt.userName, result.Name)
			require.Equal(t, 0, hasher.HashCalls)
		})
	}
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("успешное удаление", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t)
		hasher := test.NewFakePwdHasher()

		existingUser := &user.User{
			ID:           1,
			Name:         "Artem",
			Login:        "123456789",
			PasswordHash: "hash",
		}

		repo.EXPECT().
			GetByID(mock.Anything, 1).
			Return(existingUser, nil)

		repo.EXPECT().
			Delete(mock.Anything, existingUser).
			Return(nil)

		service := user.NewService(repo, hasher)

		err := service.Delete(ctx, 1)

		require.NoError(t, err)
	})

	t.Run("ошибка репозитория при поиске пользователя", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t)
		hasher := test.NewFakePwdHasher()

		expectedErr := errors.New("ошибка БД")

		repo.EXPECT().
			GetByID(mock.Anything, 1).
			Return(nil, expectedErr)

		service := user.NewService(repo, hasher)

		err := service.Delete(ctx, 1)

		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("ошибка репозитория при удалении", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t)
		hasher := test.NewFakePwdHasher()

		existingUser := &user.User{
			ID:           1,
			Name:         "Artem",
			Login:        "123456789",
			PasswordHash: "hash",
		}

		expectedErr := errors.New("ошибка БД при удалении")

		repo.EXPECT().
			GetByID(mock.Anything, 1).
			Return(existingUser, nil)

		repo.EXPECT().
			Delete(mock.Anything, existingUser).
			Return(expectedErr)

		service := user.NewService(repo, hasher)

		err := service.Delete(ctx, 1)

		require.ErrorIs(t, err, expectedErr)
	})
}

func TestService_GetByLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("пользователь найден", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t)
		hasher := test.NewFakePwdHasher()

		expectedUser := &user.User{
			ID:           1,
			Name:         "Artem",
			Login:        "login",
			PasswordHash: "hash",
		}

		repo.EXPECT().
			GetByLogin(mock.Anything, expectedUser.Login).
			Return(expectedUser, nil)

		service := user.NewService(repo, hasher)

		got, err := service.GetByLogin(ctx, "login")

		require.NoError(t, err)
		require.Equal(t, expectedUser, got)
	})

	t.Run("ошибка репозитория", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t)
		hasher := test.NewFakePwdHasher()

		expectedErr := errors.New("ошибка БД")

		repo.EXPECT().
			GetByLogin(mock.Anything, "login").
			Return(nil, expectedErr)

		service := user.NewService(repo, hasher)

		got, err := service.GetByLogin(ctx, "login")

		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, got)
	})
}

func TestService_GetByTgId(t *testing.T) {
	ctx := context.Background()

	t.Run("пользователь найден", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t)
		hasher := test.NewFakePwdHasher()

		expectedUser := &user.User{
			ID:         1,
			Name:       "Artem",
			Login:      "123456789",
			TelegramID: 123456789,
		}

		repo.EXPECT().
			GetByTgId(mock.Anything, int64(123456789)).
			Return(expectedUser, nil)

		service := user.NewService(repo, hasher)

		got, err := service.GetByTgId(ctx, 123456789)

		require.NoError(t, err)
		require.Equal(t, expectedUser, got)
	})

	t.Run("ошибка репозитория", func(t *testing.T) {
		repo := mocks.NewRepositoryMock(t)
		hasher := test.NewFakePwdHasher()

		expectedErr := errors.New("ошибка БД")

		repo.EXPECT().
			GetByTgId(mock.Anything, int64(123456789)).
			Return(nil, expectedErr)

		service := user.NewService(repo, hasher)

		got, err := service.GetByTgId(ctx, 123456789)

		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, got)
	})
}

func TestService_ValidatePassword(t *testing.T) {
	expectedErr := errors.New("ошибка проверки пароля")

	tests := []struct {
		name          string
		compareResult bool
		compareErr    error
		want          bool
		wantErr       error
	}{
		{
			name:          "пароль верный",
			compareResult: true,
			want:          true,
		},
		{
			name:          "пароль неверный",
			compareResult: false,
			want:          false,
		},
		{
			name:       "ошибка хешера",
			compareErr: expectedErr,
			wantErr:    expectedErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := test.NewFakePwdHasher()
			hasher.CompareResult = tt.compareResult
			hasher.CompareErr = tt.compareErr

			service := user.NewService(nil, hasher)

			got, err := service.ValidatePassword("hash", "password")

			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
