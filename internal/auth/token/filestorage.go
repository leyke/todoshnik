package token

import (
	"todoshnik/internal/config"
	"todoshnik/internal/storage"
)

func NewTokenStorage(c *config.Config) storage.FileStorage[Token] {
	// работа с окружением, кто сказал что эта папка/файл есть, это нужно проверять при старте приложения
	storagePath := c.App.TmpDir + "/tokens.json"
	return storage.NewFileStorage[Token](storagePath)
}
