package task

import (
	"os"
	"todoshnik/internal/storage"
)

func NewFileStorage() storage.FileStorage[Task] {
	storagePath := os.Getenv("TMP_DIR") + "/tasks.json"
	return storage.NewFileStorage[Task](storagePath)
}
