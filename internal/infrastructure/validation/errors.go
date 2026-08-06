package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	ErrNotValidate = errors.New("not validate")
)

type ValidationTag string

const (
	TagRequired ValidationTag = "required"
	TagMin      ValidationTag = "min"
	TagMax      ValidationTag = "max"
	TagNumber   ValidationTag = "number"
)

func (t ValidationTag) String() string {
	switch t {
	case TagRequired:
		return "обязательно"
	case TagMin:
		return "минимум символов"
	case TagMax:
		return "максимум символов"
	case TagNumber:
		return "должно быть числом"
	default:
		return "неверно"
	}
}

type ValidationError struct {
	err    error
	Errors map[string]string `json:"errors"`
}

func (e ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "Ошибка валидации"
	}

	var errs []error
	for field, message := range e.Errors {
		errs = append(errs, fmt.Errorf("%s: %s", field, message))
	}
	return errors.Join(errs...).Error()
}

func (e ValidationError) Unwrap() error {
	return e.err
}

func NewValidationErrorFromValidator(err error) ValidationError {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return ValidationError{
			err: ErrNotValidate,
		}
	}

	validationErrors := make(map[string]string, len(ve))
	for _, fieldErr := range ve {
		field := fieldErr.Field()

		message := getValidatorMessage(
			field,
			fieldErr.Tag(),
			fieldErr.Param(),
		)

		if existing := validationErrors[field]; existing != "" {
			validationErrors[field] = existing + "; " + message
		} else {
			validationErrors[field] = message
		}
	}

	return ValidationError{
		err:    ErrNotValidate,
		Errors: validationErrors,
	}
}

// getValidatorMessage возвращает человеко-читаемое сообщение
func getValidatorMessage(field, tag, param string) string {
	t := ValidationTag(tag)

	switch t {
	case TagMin, TagMax:
		return fmt.Sprintf("%s %s: %s", field, t, param)
	default:
		return fmt.Sprintf("%s %s", field, t)
	}
}
