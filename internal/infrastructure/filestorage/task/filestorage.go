package filestorage

import (
	"os"
	"todoshnik/internal/config"

	taskdomain "todoshnik/internal/domains/task"
	filestorage "todoshnik/internal/infrastructure/filestorage"
)

func NewFileStorage(c *config.Config) (*filestorage.FileStorage[taskdomain.Task], error) {
	storagePath := c.App.TmpDir + "/tasks.json"
	if _, err := os.Stat(storagePath); err != nil {
		return nil, err
	}

	return filestorage.NewFileStorage[taskdomain.Task](storagePath), nil
}
