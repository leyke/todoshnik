package user

import (
	"os"

	"todoshnik/internal/config"

	userdomain "todoshnik/internal/domains/user"
	filestorage "todoshnik/internal/infrastructure/filestorage"
)

func NewFileStorage(c *config.Config) (*filestorage.FileStorage[userdomain.User], error) {
	storagePath := c.App.TmpDir + "/users.json"
	if _, err := os.Stat(storagePath); err != nil {
		return nil, err
	}

	return filestorage.NewFileStorage[userdomain.User](storagePath), nil
}
