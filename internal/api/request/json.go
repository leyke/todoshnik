package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	apierrors "todoshnik/internal/api/errors"
)

func DecodeJSON[T any](r *http.Request) (*T, error) {
	var dto T

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("%w: decode json: %w", apierrors.ErrBadRequest, err)
	}

	return &dto, nil
}
