package app

import (
	"log"
	"os"
)

func NewLogger(path string) (*log.Logger, *os.File) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger, logFile
}
