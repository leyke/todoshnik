//go:build integration

package task_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TaskFixture struct {
	ID     int
	Title  string
	Done   bool
	UserID int
}

func ApplyFixtures(t *testing.T, db *sql.DB) []TaskFixture {
	t.Helper()

	var userIDs []int

	for key, login := range []string{"task_user_1", "task_user_2"} {
		var userID int

		err := db.QueryRow(`
			INSERT INTO users (name, login, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id
		`,
			"Task Test User",
			login,
			fmt.Sprintf("%s_%d", "password_hash", key),
		).Scan(&userID)

		require.NoError(t, err)

		userIDs = append(userIDs, userID)
	}

	fixtures := []TaskFixture{
		{
			Title:  "Test Task 1",
			Done:   false,
			UserID: userIDs[0],
		},
		{
			Title:  "Test Task 2",
			Done:   true,
			UserID: userIDs[0],
		},
		{
			Title:  "Test Task 3",
			Done:   false,
			UserID: userIDs[0],
		},
		{
			Title:  "Test Task 4",
			Done:   false,
			UserID: userIDs[1],
		},

		{
			Title:  "Test Task 5",
			Done:   true,
			UserID: userIDs[1],
		},
	}

	for i := range fixtures {
		err := db.QueryRow(`
			INSERT INTO tasks (title, done, user_id)
			VALUES ($1, $2, $3)
			RETURNING id
		`,
			fixtures[i].Title,
			fixtures[i].Done,
			fixtures[i].UserID,
		).Scan(&fixtures[i].ID)

		require.NoError(t, err)
	}

	return fixtures
}
