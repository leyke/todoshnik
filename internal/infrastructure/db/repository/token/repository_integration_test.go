//go:build integration

package token_test

import (
	"context"
	"testing"
	"time"

	"todoshnik/internal/infrastructure/db/repository/token"

	apptoken "todoshnik/internal/domains/token"
	tokenerror "todoshnik/internal/domains/token/errors"
	testutils "todoshnik/internal/infrastructure/utils/test"

	"github.com/stretchr/testify/require"
)

func TestDBRepository_GetByHash(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := token.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)

	tests := []struct {
		name      string
		hash      string
		wantError error
	}{
		{
			name: "exists",
			hash: fixtures[0].Hash,
		},
		{
			name:      "not found",
			hash:      "unknown",
			wantError: tokenerror.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := repo.GetByHash(ctx, tt.hash)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, token)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, token)

			require.Equal(t, tt.hash, token.Hash)
			require.Equal(t, fixtures[0].ID, token.ID)
			require.Equal(t, fixtures[0].UserID, token.UserID)
			require.Equal(t, fixtures[0].Device, token.Device)
			require.Equal(t, fixtures[0].ExpiresAt, token.ExpiresAt)
		})
	}
}

func TestDBRepository_GetByHash_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := token.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	result, err := repo.GetByHash(ctx, "unknown")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_GetAllByUserID_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := token.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)
	testUserID := fixtures[0].UserID

	var count int
	err := db.QueryRow(
		"SELECT count(*) FROM tokens WHERE user_id = $1",
		testUserID,
	).Scan(&count)

	require.NoError(t, err)
	require.Equal(t, 2, count)

	result, err := repo.GetAllByUserID(ctx, testUserID)

	require.NoError(t, err)
	require.Len(t, result, 2)

	require.Equal(t, fixtures[0].ID, result[0].ID)
	require.Equal(t, fixtures[0].UserID, result[0].UserID)
	require.Equal(t, fixtures[0].Hash, result[0].Hash)
	require.Equal(t, fixtures[0].Device, result[0].Device)
	require.Equal(t, fixtures[0].ExpiresAt, result[0].ExpiresAt)

	require.Equal(t, fixtures[1].ID, result[1].ID)
	require.Equal(t, fixtures[1].UserID, result[1].UserID)
	require.Equal(t, fixtures[1].Hash, result[1].Hash)
	require.Equal(t, fixtures[1].Device, result[1].Device)
	require.Equal(t, fixtures[1].ExpiresAt, result[1].ExpiresAt)
}

func TestDBRepository_GetAllByUserID_Empty(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := token.NewRepository(db)
	ctx := context.Background()
	testUserID := 1

	result, err := repo.GetAllByUserID(ctx, testUserID)

	require.NoError(t, err)
	require.Empty(t, result)
}

func TestDBRepository_GetAllByUserID_DbError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := token.NewRepository(db)
	ctx := context.Background()
	testUserID := 1

	require.NoError(t, db.Close())

	result, err := repo.GetAllByUserID(ctx, testUserID)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_GetExpiredTokens_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := token.NewRepository(db)
	ctx := context.Background()
	now := time.Now().Round(0)
	_ = ApplyFixtures(t, db)

	result, err := repo.GetExpiredTokens(ctx, now)

	require.NoError(t, err)

	for _, token := range result {
		require.Less(t, token.ExpiresAt, now)
	}
}

func TestDBRepository_GetExpiredTokens_Empty(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := token.NewRepository(db)
	ctx := context.Background()
	now := time.Now().Round(0)

	result, err := repo.GetExpiredTokens(ctx, now)

	require.NoError(t, err)
	require.Empty(t, result)
}

func TestDBRepository_GetExpiredTokens_DbError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := token.NewRepository(db)
	ctx := context.Background()
	now := time.Now().Round(0)

	require.NoError(t, db.Close())

	result, err := repo.GetExpiredTokens(ctx, now)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_Delete_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := token.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)

	input := &apptoken.Token{
		ID: fixtures[0].ID,
	}

	err := repo.Delete(ctx, input)
	require.NoError(t, err)

	result, err := repo.GetByHash(ctx, input.Hash)

	require.ErrorIs(t, err, tokenerror.ErrNotFound)
	require.Nil(t, result)
}

func TestDBRepository_Delete_NotFound(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := token.NewRepository(db)
	ctx := context.Background()

	input := &apptoken.Token{
		ID: 99999,
	}

	err := repo.Delete(ctx, input)

	require.ErrorIs(t, err, tokenerror.ErrNotFound)
}

func TestDBRepository_Delete_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := token.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	input := &apptoken.Token{
		ID: 99999,
	}

	err := repo.Delete(ctx, input)

	require.Error(t, err)
}

func TestDBRepository_Create_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := token.NewRepository(db)
	ctx := context.Background()
	now := time.Now().Round(0)

	fixtures := ApplyFixtures(t, db)

	input := &apptoken.Token{
		UserID:    fixtures[0].UserID,
		Hash:      "new_hash",
		Device:    apptoken.DeviceTypeBot,
		ExpiresAt: now,
	}

	result, err := repo.Create(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.NotZero(t, result.ID)
	require.Equal(t, input.Hash, result.Hash)
	require.Equal(t, input.UserID, result.UserID)
	require.Equal(t, input.Device, result.Device)
	require.Equal(t, input.ExpiresAt, result.ExpiresAt)
}

func TestDBRepository_Create_UserNotFound(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := token.NewRepository(db)
	ctx := context.Background()
	now := time.Now().Round(0)

	input := &apptoken.Token{
		UserID:    99999,
		Hash:      "new_hash",
		Device:    apptoken.DeviceTypeApi,
		ExpiresAt: now,
	}

	result, err := repo.Create(ctx, input)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_Create_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := token.NewRepository(db)
	ctx := context.Background()
	now := time.Now().Round(0)

	require.NoError(t, db.Close())

	input := &apptoken.Token{
		UserID:    99999,
		Hash:      "new_hash",
		Device:    apptoken.DeviceTypeApi,
		ExpiresAt: now,
	}

	result, err := repo.Create(ctx, input)

	require.Error(t, err)
	require.Nil(t, result)
}
