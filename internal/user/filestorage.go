package user

import (
	"todoshnik/internal/config"
	"todoshnik/internal/storage"
)

func NewFileStorage(c *config.Config) storage.FileStorage[User] {
	storagePath := c.App.TmpDir + "/users.json"
	return storage.NewFileStorage[User](storagePath)
}
