package validator

import (
	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	v *validator.Validate
}

func NewStructValidator() *StructValidator {
	return &StructValidator{
		v: validator.New(),
	}
}

func (sv *StructValidator) Validate(s interface{}) error {
	err := sv.v.Struct(s)

	switch err.(type) {
	case *validator.InvalidValidationError:
		panic(err)
	case validator.ValidationErrors:
		return err
	}

	return nil
}
