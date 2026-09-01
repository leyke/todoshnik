package token

import (
	"os"
	"todoshnik/internal/config"

	tokendomain "todoshnik/internal/domains/token"
	filestorage "todoshnik/internal/infrastructure/filestorage"
)

func NewTokenStorage(c *config.Config) (*filestorage.FileStorage[tokendomain.Token], error) {
	storagePath := c.App.TmpDir + "/tokens.json"
	if _, err := os.Stat(storagePath); err != nil {
		return nil, err
	}

	return filestorage.NewFileStorage[tokendomain.Token](storagePath), nil
}
