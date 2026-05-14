package request

import (
	"encoding/json"
	"net/http"
	apperrors "todoshnik/internal/errors"
)

func DecodeJSON[T any](r *http.Request) (*T, error) {
	var dto T

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		if err.Error() == "EOF" {
			return nil, apperrors.ErrEmptyBody
		}

		return nil, apperrors.ErrInvalidJSON
	}

	return &dto, nil
}
