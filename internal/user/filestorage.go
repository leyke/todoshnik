package user

import (
	"os"
	"todoshnik/internal/storage"
)

func NewFileStorage() storage.FileStorage[User] {
	storagePath := os.Getenv("TMP_DIR") + "/users.json"
	return storage.NewFileStorage[User](storagePath)
}
