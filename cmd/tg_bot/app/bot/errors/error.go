package errors

import "errors"

var (
	ErrUnknownMethod = errors.New("unknown method")
	ErrInvalidTaskID = errors.New("invalid task id")
)
