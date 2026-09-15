package app

import (
	"log"
	"os"
)

func NewLogger() *log.Logger {
	logger := log.New(
		os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	return logger
}
