package validation

import (
	"github.com/go-playground/validator/v10"
)

var valid = *validator.New()

func Validate(v any) error {
	return valid.Struct(v)
}
