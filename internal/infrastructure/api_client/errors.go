package client

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
	ErrUnAuth   = errors.New("unauthorized")
)

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("api error: %d: %s", e.StatusCode, e.Message)
}
