package common

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	sv *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		sv: validator.New(),
	}
}

func (v *Validator) Validate(s interface{}) error {
	err := v.sv.Struct(s)

	switch err.(type) {
	case *validator.InvalidValidationError:
		panic(err)
	case validator.ValidationErrors:
		return err
	}

	return nil
}
