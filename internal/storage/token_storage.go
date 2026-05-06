package storage

import (
	"os"
	"todoshnik/internal/domain"
)

func NewTokenStorage() FileStorage[domain.Token] {
	storagePath := os.Getenv("TMP_DIR") + "/tokens.json"
	return NewFileStorage[domain.Token](storagePath)
}
