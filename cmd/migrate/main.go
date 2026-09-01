package main

import (
	"fmt"
	"log"

	"todoshnik/internal/config"
	"todoshnik/internal/infrastructure/db"

	"github.com/pressly/goose/v3"
)

var (
	AppVersion = "dev"
	CommitHash = "unknown"
)

func main() {
	fmt.Println("Version:", AppVersion)
	fmt.Println("Commit :", CommitHash)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDb, err := db.NewGormDb(cfg)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := gormDb.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(
		conn,
		"internal/infrastructure/db/migrations",
	); err != nil {
		log.Fatal(err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations applied")
}
