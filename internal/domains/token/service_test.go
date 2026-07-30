package token

import (
	"context"
	"errors"
	"testing"
	"time"

	userdomain "todoshnik/internal/domains/user"
	utilstest "todoshnik/internal/infrastructure/utils/test"

	"github.com/stretchr/testify/require"
)

const (
	testUserID         = 10
	testRawToken       = "raw-token"
	testTokenHash      = "hashed-token"
	testErrorTokenHash = "bad-token"

	testDevice = DeviceType("web")
)

type repositoryMock struct {
	createFunc           func(ctx context.Context, token *Token) (*Token, error)
	getFunc              func(ctx context.Context, hash string) (*Token, error)
	getExpiredTokensFunc func(ctx context.Context, before time.Time) ([]*Token, error)
	deleteFunc           func(ctx context.Context, token *Token) error
}

func (m *repositoryMock) GetAllByUserID(ctx context.Context, id int) ([]*Token, error) {
	return nil, nil
}

func (m *repositoryMock) GetByHash(ctx context.Context, hash string) (*Token, error) {
	if m.getFunc == nil {
		return nil, nil
	}

	return m.getFunc(ctx, hash)
}

func (m *repositoryMock) GetExpiredTokens(ctx context.Context, before time.Time) ([]*Token, error) {
	if m.getExpiredTokensFunc == nil {
		return nil, nil
	}

	return m.getExpiredTokensFunc(ctx, before)
}

func (m *repositoryMock) Create(ctx context.Context, token *Token) (*Token, error) {
	if m.createFunc == nil {
		return nil, nil
	}

	return m.createFunc(ctx, token)
}

func (m *repositoryMock) Delete(ctx context.Context, token *Token) error {
	if m.deleteFunc == nil {
		return nil
	}

	return m.deleteFunc(ctx, token)
}

func TestService_Add(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	now := time.Date(2026, 7, 29, 10, 0, 0, 0, time.UTC)

	user := &userdomain.User{
		ID: testUserID,
	}

	t.Run("Успешное создание токена", func(t *testing.T) {
		t.Parallel()

		repo := &repositoryMock{
			createFunc: func(ctx context.Context, token *Token) (*Token, error) {
				require.Equal(t, testUserID, token.UserID)
				require.Equal(t, testTokenHash, token.Hash)
				require.Equal(t, testDevice, token.Device)
				require.Equal(t,
					now.Add(time.Hour),
					token.ExpiresAt,
				)

				return token, nil
			},
		}

		service := NewService(
			repo,
			&utilstest.FakeTokenGenerator{
				RawToken:    testRawToken,
				HashedToken: testTokenHash,
			},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		result, err := service.Add(
			ctx,
			user,
			testDevice,
		)

		require.NoError(t, err)
		require.Equal(t, testRawToken, result)
	})

	t.Run("Ошибка репозитория", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("Ошибка БД")

		service := NewService(
			&repositoryMock{
				createFunc: func(ctx context.Context, token *Token) (*Token, error) {
					return nil, expectedErr
				},
			},
			&utilstest.FakeTokenGenerator{},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		_, err := service.Add(
			ctx,
			user,
			testDevice,
		)

		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("Ошибка генерации токена", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("Ошибка генерации")

		repo := &repositoryMock{
			createFunc: func(ctx context.Context, token *Token) (*Token, error) {
				t.Fatal("repository Create should not be called")
				return nil, nil
			},
		}

		service := NewService(
			repo,
			&utilstest.FakeTokenGenerator{
				GenerateErr: expectedErr,
			},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		_, err := service.Add(
			ctx,
			user,
			testDevice,
		)

		require.ErrorIs(t, err, expectedErr)
	})
}

func TestService_Get(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	now := time.Date(2026, 7, 29, 10, 0, 0, 0, time.UTC)

	expected := &Token{
		UserID:    testUserID,
		Hash:      testTokenHash,
		Device:    testDevice,
		ExpiresAt: now.Add(time.Hour),
	}

	t.Run("Успешное получение токена", func(t *testing.T) {
		t.Parallel()

		repo := &repositoryMock{
			getFunc: func(ctx context.Context, hash string) (*Token, error) {
				require.Equal(t, testTokenHash, hash)
				return &Token{
					UserID:    testUserID,
					Hash:      testTokenHash,
					Device:    testDevice,
					ExpiresAt: now.Add(time.Hour),
				}, nil
			},
		}

		service := NewService(
			repo,
			&utilstest.FakeTokenGenerator{
				RawToken:    testRawToken,
				HashedToken: testTokenHash,
			},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		actual, err := service.Get(
			ctx,
			testRawToken,
		)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("Ошибка репозитория", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("Ошибка БД")

		service := NewService(
			&repositoryMock{
				getFunc: func(ctx context.Context, hash string) (*Token, error) {
					require.Equal(t, testTokenHash, hash)
					return nil, expectedErr
				},
			},
			&utilstest.FakeTokenGenerator{
				HashedToken: testTokenHash,
			},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		_, err := service.Get(
			ctx,
			testRawToken,
		)

		require.ErrorIs(t, err, expectedErr)
	})
}

func TestService_ClearExpiredTokens(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	now := time.Date(2026, 7, 29, 10, 0, 0, 0, time.UTC)

	mockedToken := &Token{
		UserID:    testUserID,
		Hash:      testTokenHash,
		Device:    testDevice,
		ExpiresAt: now.Add(-time.Hour),
	}

	mockedErrorToken := &Token{
		UserID:    testUserID,
		Hash:      testErrorTokenHash,
		Device:    testDevice,
		ExpiresAt: now.Add(-time.Hour),
	}

	t.Run("Несколько токенов, все удалены", func(t *testing.T) {
		t.Parallel()

		expectCount := 3

		repo := &repositoryMock{
			getExpiredTokensFunc: func(ctx context.Context, before time.Time) ([]*Token, error) {
				require.Equal(t, now, before)

				items := make([]*Token, 0, expectCount)
				for range expectCount {
					items = append(items, mockedToken)
				}
				return items, nil
			},
		}

		service := NewService(
			repo,
			&utilstest.FakeTokenGenerator{
				RawToken:    testRawToken,
				HashedToken: testTokenHash,
			},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		actual, err := service.ClearExpiredTokens(ctx)

		require.NoError(t, err)
		require.Equal(t, expectCount, actual)
	})

	t.Run("Один просроченный токен, удаление успешно", func(t *testing.T) {
		t.Parallel()

		expectCount := 1

		repo := &repositoryMock{
			getExpiredTokensFunc: func(ctx context.Context, before time.Time) ([]*Token, error) {
				require.Equal(t, now, before)

				items := make([]*Token, 0, expectCount)
				for range expectCount {
					items = append(items, mockedToken)
				}
				return items, nil
			},
		}

		service := NewService(
			repo,
			&utilstest.FakeTokenGenerator{
				RawToken:    testRawToken,
				HashedToken: testTokenHash,
			},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		actual, err := service.ClearExpiredTokens(ctx)

		require.NoError(t, err)
		require.Equal(t, expectCount, actual)
	})

	t.Run("Нет просроченных токенов", func(t *testing.T) {
		t.Parallel()
		repo := &repositoryMock{
			getExpiredTokensFunc: func(ctx context.Context, before time.Time) ([]*Token, error) {
				require.Equal(t, now, before)

				return nil, nil
			},
		}

		service := NewService(
			repo,
			&utilstest.FakeTokenGenerator{
				RawToken:    testRawToken,
				HashedToken: testTokenHash,
			},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		actual, err := service.ClearExpiredTokens(ctx)

		require.NoError(t, err)
		require.Equal(t, 0, actual)
	})

	t.Run("Ошибка репозитория", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("Ошибка БД")

		service := NewService(
			&repositoryMock{
				getExpiredTokensFunc: func(ctx context.Context, before time.Time) ([]*Token, error) {
					require.Equal(t, now, before)

					return nil, expectedErr
				},
			},
			&utilstest.FakeTokenGenerator{},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		_, err := service.ClearExpiredTokens(ctx)

		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("Несколько токенов, 1 ошибка", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("Ошибка БД")
		expectedErrCount := 1
		expectSuccessCount := 2

		service := NewService(
			&repositoryMock{
				getExpiredTokensFunc: func(ctx context.Context, before time.Time) ([]*Token, error) {
					require.Equal(t, now, before)

					items := make([]*Token, 0, expectSuccessCount+expectedErrCount)
					for range expectSuccessCount {
						items = append(items, mockedToken)
					}
					for range expectedErrCount {
						items = append(items, mockedErrorToken)
					}
					return items, nil
				},
				deleteFunc: func(ctx context.Context, token *Token) error {
					if token.Hash == testErrorTokenHash {
						return expectedErr
					}
					return nil
				},
			},
			&utilstest.FakeTokenGenerator{},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		deletedCount, err := service.ClearExpiredTokens(ctx)

		require.Equal(t, expectSuccessCount, deletedCount)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("Несколько токенов, все удаления завершились ошибкой", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("Ошибка БД")
		expectedErrCount := 3
		expectedSuccessCount := 0

		service := NewService(
			&repositoryMock{
				getExpiredTokensFunc: func(ctx context.Context, before time.Time) ([]*Token, error) {
					require.Equal(t, now, before)

					items := make([]*Token, 0, expectedErrCount)
					for range expectedErrCount {
						items = append(items, mockedErrorToken)
					}
					return items, nil
				},
				deleteFunc: func(ctx context.Context, token *Token) error {
					if token.Hash == testErrorTokenHash {
						return expectedErr
					}
					return nil
				},
			},
			&utilstest.FakeTokenGenerator{},
			&utilstest.FakeClock{Time: now},
			Config{
				Ttl: time.Hour,
			},
		)

		deletedCount, err := service.ClearExpiredTokens(ctx)

		require.Equal(t, expectedSuccessCount, deletedCount)
		require.ErrorIs(t, err, expectedErr)
	})
}
