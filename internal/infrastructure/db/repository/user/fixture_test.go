//go:build integration

package user_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserFixture struct {
	ID           int
	Name         string
	Login        string
	PasswordHash string
	TelegramID   int64
}

func ApplyFixtures(t *testing.T, db *sql.DB) []UserFixture {
	t.Helper()

	fixtures := []UserFixture{
		{
			Name:         "Test User 1",
			Login:        "test_user_1",
			PasswordHash: "password_hash_1",
			TelegramID:   111111111,
		},
		{
			Name:         "Test User 2",
			Login:        "test_user_2",
			PasswordHash: "password_hash_2",
			TelegramID:   222222222,
		},
	}

	for i := range fixtures {
		err := db.QueryRow(`
			INSERT INTO users (name, login, password_hash, telegram_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`,
			fixtures[i].Name,
			fixtures[i].Login,
			fixtures[i].PasswordHash,
			fixtures[i].TelegramID,
		).Scan(&fixtures[i].ID)

		require.NoError(t, err)
	}

	return fixtures
}
