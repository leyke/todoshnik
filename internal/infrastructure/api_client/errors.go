package client

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrUnAuth   = errors.New("unauthorized")
)
