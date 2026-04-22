package storage

type Storage interface {
	Save(any) error
	Load() (any, error)
}
