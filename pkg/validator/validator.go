package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

func New() *Validator {
	return &Validator{
		validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v Validator) ValidateStruct(s interface{}) error {
	err := v.Struct(s)

	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)

	if len(validationErrors) < 1 {
		return nil
	}

	return MessageForFieldError(validationErrors[0])
}

func MessageForFieldError(e validator.FieldError) error {
	switch e.Tag() {
	case "required":
		return fmt.Errorf("field '%s' : is required", e.Field())
	case "email":
		return fmt.Errorf("field '%s' : must be valid email", e.Field())
	case "min":
		return fmt.Errorf("field '%s' : Have min length of %s", e.Field(), e.Param())
	case "max":
		return fmt.Errorf("field '%s' : Have max length of %s", e.Field(), e.Param())
	}
	return e
}
