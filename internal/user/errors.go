package user

import (
	"errors"
)

var (
	ErrConflict = errors.New("conflict")
	ErrNotFound = errors.New("user not found")
)
