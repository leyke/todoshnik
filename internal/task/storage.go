package task

type TaskStorage interface {
	Save(tasks map[int]*Task) error
	Load() (map[int]*Task, error)
}
