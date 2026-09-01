package filestorage

type Storage interface {
	Save(any) error
	Load() (any, error)
}
