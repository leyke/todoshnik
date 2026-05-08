package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDb() (*gorm.DB, error) {
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}
