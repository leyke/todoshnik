package task

import (
	"errors"
)

var (
	ErrInvalidTaskID = errors.New("invalid task id")
	ErrNotFound      = errors.New("task not found")
)
