package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Struct[T any](field T) error {
	if err := validate.Struct(field); err != nil {
		return err
	}
	return nil
}
