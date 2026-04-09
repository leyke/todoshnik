package storage

import "todoshnik/internal/domain"

type TaskStorage interface {
	Save(tasks map[int]*domain.Task) error
	Load() (map[int]*domain.Task, error)
}
