package test_test

import (
	"testing"
	"todoshnik/internal/infrastructure/utils/test"

	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	db := test.TestDB(t)

	var result int

	err := db.QueryRow("SELECT 1").Scan(&result)

	require.NoError(t, err)
	require.Equal(t, 1, result)
}
