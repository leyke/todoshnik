package storage

import (
	"encoding/json"
	"os"
	"todoshnik/internal/domain"
)

type FileStorage struct {
	filename string
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{
		filename: filename,
	}
}

func (fs *FileStorage) Save(tasks map[int]*domain.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.filename, data, 0644)
}

func (fs *FileStorage) Load() (map[int]*domain.Task, error) {
	data, err := os.ReadFile(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[int]*domain.Task), nil
		}
		return nil, err
	}

	tasks := make(map[int]*domain.Task)
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}
