package task

import "todoshnik/internal/infrastructure/identity"

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
)

type TaskFilter struct {
	Status Status
	Scope  identity.AccessScope
}
