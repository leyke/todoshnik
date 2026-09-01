package errors

import (
	"errors"
)

var (
	ErrConflict = errors.New("conflict")
	ErrNotFound = errors.New("user not found")
)
