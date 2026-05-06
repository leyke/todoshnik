package storage

import (
	"os"
	"todoshnik/internal/domain"
)

func NewUserStorage() FileStorage[domain.User] {
	storagePath := os.Getenv("TMP_DIR") + "/users.json"
	return NewFileStorage[domain.User](storagePath)
}
