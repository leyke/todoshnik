package task

import (
	"encoding/json"
	"os"
)

type FileStorage struct {
	filename string
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{
		filename: filename,
	}
}

func (fs *FileStorage) Save(tasks map[int]*Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.filename, data, 0644)
}

func (fs *FileStorage) Load() (map[int]*Task, error) {
	data, err := os.ReadFile(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[int]*Task), nil
		}
		return nil, err
	}

	tasks := make(map[int]*Task)
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}
