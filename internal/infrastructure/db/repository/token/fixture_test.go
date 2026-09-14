//go:build integration

package token_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"todoshnik/internal/domains/token"

	"github.com/stretchr/testify/require"
)

type TokenFixture struct {
	ID        int
	UserID    int
	Hash      string
	Device    token.DeviceType
	ExpiresAt time.Time
}

func ApplyFixtures(t *testing.T, db *sql.DB) []TokenFixture {
	t.Helper()

	var userIDs []int

	for key, login := range []string{"token_user_1", "token_user_2"} {
		var userID int

		err := db.QueryRow(`
			INSERT INTO users (name, login, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id
		`,
			"Token Test User",
			login,
			fmt.Sprintf("%s_%d", "password_hash", key),
		).Scan(&userID)

		require.NoError(t, err)

		userIDs = append(userIDs, userID)
	}

	expiresAt := time.Now().Add(24 * time.Hour).Round(0)
	expiredTime := time.Now().Add(-24 * time.Hour).Round(0)

	fixtures := []TokenFixture{
		{
			UserID:    userIDs[0],
			Hash:      "token_hash_1",
			Device:    token.DeviceTypeApi,
			ExpiresAt: expiresAt,
		},
		{
			UserID:    userIDs[0],
			Hash:      "token_hash_2",
			Device:    token.DeviceTypeBot,
			ExpiresAt: expiresAt,
		},
		{
			UserID:    userIDs[1],
			Hash:      "token_hash_3",
			Device:    token.DeviceTypeBot,
			ExpiresAt: expiredTime,
		},
	}

	for i := range fixtures {
		err := db.QueryRow(`
			INSERT INTO tokens (user_id, hash, device, expires_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`,
			fixtures[i].UserID,
			fixtures[i].Hash,
			fixtures[i].Device,
			fixtures[i].ExpiresAt,
		).Scan(&fixtures[i].ID)

		require.NoError(t, err)
	}

	return fixtures
}
