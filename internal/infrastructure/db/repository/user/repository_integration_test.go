//go:build integration

package user_test

import (
	"context"
	"testing"

	"todoshnik/internal/infrastructure/db/repository/user"

	appuser "todoshnik/internal/domains/user"
	usererrors "todoshnik/internal/domains/user/errors"
	testutils "todoshnik/internal/infrastructure/utils/test"

	"github.com/stretchr/testify/require"
)

func TestDBRepository_GetByID(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	users := ApplyFixtures(t, db)

	tests := []struct {
		name      string
		id        int
		wantError error
	}{
		{
			name: "user exists",
			id:   users[0].ID,
		},
		{
			name:      "user not found",
			id:        999999,
			wantError: usererrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(ctx, tt.id)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)

			require.Equal(t, tt.id, user.ID)
			require.Equal(t, users[0].Name, user.Name)
			require.Equal(t, users[0].Login, user.Login)
			require.Equal(t, users[0].TelegramID, user.TelegramID)
		})
	}
}

func TestDBRepository_GetByID_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	result, err := repo.GetByID(ctx, 1)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_GetByLogin(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	users := ApplyFixtures(t, db)

	tests := []struct {
		name      string
		login     string
		wantError error
	}{
		{
			name:  "user exists",
			login: users[0].Login,
		},
		{
			name:      "user not found",
			login:     "testtttt",
			wantError: usererrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByLogin(ctx, tt.login)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, users[0].ID, user.ID)
			require.Equal(t, users[0].Name, user.Name)
			require.Equal(t, users[0].Login, user.Login)
			require.Equal(t, users[0].TelegramID, user.TelegramID)
		})
	}
}

func TestDBRepository_GetByLogin_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	result, err := repo.GetByLogin(ctx, "UnknownLogin")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_GetByTgId(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	users := ApplyFixtures(t, db)

	tests := []struct {
		name       string
		telegramID int64
		wantError  error
	}{
		{
			name:       "user exists",
			telegramID: users[0].TelegramID,
		},
		{
			name:       "user not found",
			telegramID: 3,
			wantError:  usererrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByTgId(ctx, tt.telegramID)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, users[0].ID, user.ID)
			require.Equal(t, users[0].Name, user.Name)
			require.Equal(t, users[0].Login, user.Login)
			require.Equal(t, users[0].TelegramID, user.TelegramID)
		})
	}
}

func TestDBRepository_GetByTgId_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	result, err := repo.GetByTgId(ctx, 1)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_Create_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	input := &appuser.User{
		Name:         "New User",
		Login:        "new_user",
		PasswordHash: "password_hash",
		TelegramID:   123456789,
	}

	result, err := repo.Create(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.NotZero(t, result.ID)
	require.Equal(t, input.Name, result.Name)
	require.Equal(t, input.Login, result.Login)
	require.Equal(t, input.PasswordHash, result.PasswordHash)
	require.Equal(t, input.TelegramID, result.TelegramID)
}

func TestDBRepository_Create_DuplicateLogin(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	users := ApplyFixtures(t, db)

	input := &appuser.User{
		Name:         "Another User",
		Login:        users[0].Login,
		PasswordHash: "another_password",
		TelegramID:   999999999,
	}

	result, err := repo.Create(ctx, input)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_Create_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	input := &appuser.User{
		Name:         "New User",
		Login:        "new_user",
		PasswordHash: "password_hash",
		TelegramID:   123456789,
	}

	result, err := repo.Create(ctx, input)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_Update_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	users := ApplyFixtures(t, db)

	input := &appuser.User{
		ID:           users[0].ID,
		Name:         "Updated User",
		Login:        "updated_login",
		PasswordHash: "updated_password",
		TelegramID:   999999999,
	}

	err := repo.Update(ctx, input)
	require.NoError(t, err)

	result, err := repo.GetByID(ctx, input.ID)

	require.NoError(t, err)
	require.Equal(t, input.Name, result.Name)
	require.Equal(t, input.Login, result.Login)
	require.Equal(t, input.TelegramID, result.TelegramID)
}

func TestDBRepository_Update_NotFound(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	input := &appuser.User{
		ID:           999999,
		Name:         "Unknown User",
		Login:        "unknown_user",
		PasswordHash: "password_hash",
		TelegramID:   123456789,
	}

	err := repo.Update(ctx, input)

	require.ErrorIs(t, err, usererrors.ErrNotFound)
}

func TestDBRepository_Update_DuplicateLogin(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	users := ApplyFixtures(t, db)

	input := &appuser.User{
		Name:         "Another User",
		Login:        users[0].Login,
		PasswordHash: "another_password",
		TelegramID:   999999999,
	}

	err := repo.Update(ctx, input)

	require.Error(t, err)
}

func TestDBRepository_Update_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	input := &appuser.User{
		ID:           1,
		Name:         "Updated User",
		Login:        "updated_user",
		PasswordHash: "password_hash",
		TelegramID:   123456789,
	}

	err := repo.Update(ctx, input)

	require.Error(t, err)
}

func TestDBRepository_Delete_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	users := ApplyFixtures(t, db)

	input := &appuser.User{
		ID: users[0].ID,
	}

	err := repo.Delete(ctx, input)
	require.NoError(t, err)

	result, err := repo.GetByID(ctx, input.ID)

	require.ErrorIs(t, err, usererrors.ErrNotFound)
	require.Nil(t, result)
}

func TestDBRepository_Delete_NotFound(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	input := &appuser.User{
		ID: 999999,
	}

	err := repo.Delete(ctx, input)

	require.ErrorIs(t, err, usererrors.ErrNotFound)
}

func TestDBRepository_Delete_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	input := &appuser.User{
		ID: 1,
	}

	err := repo.Delete(ctx, input)

	require.Error(t, err)
}

func TestDBRepository_List_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)

	result, err := repo.List(ctx)

	require.NoError(t, err)
	require.Len(t, result, 2)

	require.Equal(t, fixtures[0].ID, result[0].ID)
	require.Equal(t, fixtures[0].Name, result[0].Name)
	require.Equal(t, fixtures[0].Login, result[0].Login)
	require.Equal(t, fixtures[0].TelegramID, result[0].TelegramID)

	require.Equal(t, fixtures[1].ID, result[1].ID)
	require.Equal(t, fixtures[1].Name, result[1].Name)
	require.Equal(t, fixtures[1].Login, result[1].Login)
	require.Equal(t, fixtures[1].TelegramID, result[1].TelegramID)
}

func TestDBRepository_List_Empty(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	result, err := repo.List(ctx)

	require.NoError(t, err)
	require.Empty(t, result)
}

func TestDBRepository_List_DbError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := user.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	result, err := repo.List(ctx)

	require.Error(t, err)
	require.Nil(t, result)
}
