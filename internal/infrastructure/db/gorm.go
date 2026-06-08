package db

import (
	"todoshnik/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDb(cfg *config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.Postgres.Host + " user=" + cfg.Postgres.User + " password=" + cfg.Postgres.Password + " dbname=" + cfg.Postgres.DBName + " port=" + cfg.Postgres.Port + " sslmode=" + cfg.Postgres.SSLMode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}
