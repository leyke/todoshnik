package user

import (
	"todoshnik/internal/config"
	"todoshnik/internal/storage"
)

func NewFileStorage(c *config.Config) storage.FileStorage[User] {
	// TODO чекнуть наличие файла

	storagePath := c.App.TmpDir + "/users.json"
	return storage.NewFileStorage[User](storagePath)
}
