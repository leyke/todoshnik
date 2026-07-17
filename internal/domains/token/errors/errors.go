package errors

import "errors"

var (
	ErrUserNotFound = errors.New("user for token not found")
	ErrNotFound     = errors.New("token not found")
)
