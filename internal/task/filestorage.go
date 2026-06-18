package task

import (
	"todoshnik/internal/config"
	"todoshnik/internal/storage"
)

func NewFileStorage(c *config.Config) storage.FileStorage[Task] {
	// TODO чекнуть наличие файла

	storagePath := c.App.TmpDir + "/tasks.json"
	return storage.NewFileStorage[Task](storagePath)
}
