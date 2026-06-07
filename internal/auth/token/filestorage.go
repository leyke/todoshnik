package token

import (
	"os"
	"todoshnik/internal/storage"
)

func NewTokenStorage() storage.FileStorage[Token] {
	// работа с окружением, кто сказал что эта папка/файл есть, это нужно проверять при старте приложения
	storagePath := os.Getenv("TMP_DIR") + "/tokens.json"
	return storage.NewFileStorage[Token](storagePath)
}
