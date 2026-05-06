package storage

import (
	"os"
	"todoshnik/internal/domain"
)

func NewTaskStorage() FileStorage[domain.Task] {
	storagePath := os.Getenv("TMP_DIR") + "/tasks.json"
	return NewFileStorage[domain.Task](storagePath)
}
