package db

import (
	"database/sql"
	"fmt"
	"todoshnik/internal/config"

	_ "github.com/jackc/pgx/v5" // драйвер для соединения с бд
)

func New(cfg config.Config) (*sql.DB, error) {
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
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Postgres.DbMaxOpenConnections)
	db.SetMaxIdleConns(cfg.Postgres.DbMaxIdleConnections)

	return db, nil
}
