package storage

import (
	"encoding/json"
	"os"
)

type FileStorage[T any] struct {
	filename string
}

func NewFileStorage[T any](filename string) FileStorage[T] {
	return FileStorage[T]{
		filename: filename,
	}
}

func (fs *FileStorage[T]) Save(data map[int]*T) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.filename, jsonData, 0644)
}

func (fs *FileStorage[T]) Load() (map[int]*T, error) {
	data, err := os.ReadFile(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[int]*T), nil
		}
		return nil, err
	}

	result := make(map[int]*T)
	err = json.Unmarshal(data, &result)
	return result, err
}
