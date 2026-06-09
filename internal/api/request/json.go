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
		// 1. разве тут может быть EOF? Стоит написать тесты чтобы проверить, что-то мне кажется что не может и даже
		// пустая строка будет валидным json, может быть невалидный из одно открывающей/закрывающей фигурной скобки,
		// но там будет анекспектед энд оф json. Возможно я ошибаюсь и давно не работал со стандартным json маршалером.
		// В любом случае я бы перепроверил тестом.
		// 2. Почму просто не заврапать ошибку return nil, fmt.Printf("DecodeJSON error: %w", err)
		if err.Error() == "EOF" {
			return nil, apperrors.ErrEmptyBody
		}

		return nil, apperrors.ErrInvalidJSON
	}

	return &dto, nil
}
