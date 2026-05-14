package token

import (
	"os"
	"todoshnik/internal/storage"
)

func NewTokenStorage() storage.FileStorage[Token] {
	storagePath := os.Getenv("TMP_DIR") + "/tokens.json"
	return storage.NewFileStorage[Token](storagePath)
}
