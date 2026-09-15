package test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"todoshnik/internal/config"

	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) *sql.DB {
	t.Helper()

	require.NoError(t, loadEnv())

	cfg, err := config.Load()
	require.NoError(t, err)

	cfg.Postgres.DBName = "tododb_test"
	cfg.Postgres.Host = "localhost"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.Port,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)

	require.NoError(t, db.Ping())

	return db
}

func TestDBWithCleanup(t *testing.T) *sql.DB {
	t.Helper()

	db := TestDB(t)

	t.Cleanup(func() {
		Cleanup(t, db)
	})

	return db
}

func Cleanup(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
		TRUNCATE users RESTART IDENTITY CASCADE;
		TRUNCATE tasks RESTART IDENTITY CASCADE;
		TRUNCATE tokens RESTART IDENTITY CASCADE;
	`)

	require.NoError(t, err)
}
