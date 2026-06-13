package errors

import "errors"

var (
	ErrBadRequest  = errors.New("bad request")
	ErrNotFound    = errors.New("not found")
	ErrEmptyBody   = errors.New("empty body")
	ErrInvalidJSON = errors.New("invalid json")
	ErrUnauth      = errors.New("unauthorized")
	ErrConflict    = errors.New("conflict")
)
