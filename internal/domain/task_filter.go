package domain

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
)

type TaskFilter struct {
	Status TaskStatus
	Scope  AccessScope
}
