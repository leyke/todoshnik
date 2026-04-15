package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrNotValidate = errors.New("not validate")
)

type ValidationError struct {
	error
	Errors map[string]string `json:"errors"`
}

func (e ValidationError) Is(target error) bool {
	return target == ErrNotValidate
}

func (e ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "Ошибка валидации"
	}

	var errs []string
	for field, message := range e.Errors {
		errs = append(errs, fmt.Sprintf("%s: %s", field, message))
	}
	return "Ошибка валидации: " + strings.Join(errs, "; ")

}
func NewValidationErrorFromValidator(ve validator.ValidationErrors) ValidationError {
	errors := make(map[string]string)

	for _, err := range ve {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		message := getValidatorMessage(field, tag, param)
		errors[field] = message
	}
	fmt.Println(errors)
	return ValidationError{
		error:  ErrNotValidate,
		Errors: errors,
	}
}

// getValidatorMessage возвращает человеко-читаемое сообщение
func getValidatorMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s обязательно", field)
	case "min":
		return fmt.Sprintf("%s минимум символов: %s", field, param)
	case "max":
		return fmt.Sprintf("%s максимум символов: %s", field, param)
	case "number":
		return fmt.Sprintf("%s должно быть числом", field)
	default:
		return fmt.Sprintf("%s неверно (%s)", field, tag)
	}
}
