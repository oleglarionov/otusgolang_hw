package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"strings"
)

type InStringValidator struct {
	values []string
}

func (v *InStringValidator) Type() ValidatorType {
	return StringValidator
}

func (v *InStringValidator) TagName() string {
	return "in"
}

func (v *InStringValidator) Build(constraint string) {
	v.values = strings.Split(constraint, ",")
}

func (v *InStringValidator) Validate(value interface{}) error {
	casted := value.(string)
	for _, value := range v.values {
		if casted == value {
			return nil
		}
	}

	return InStringValidatorError{
		Required: v.values,
		Actual:   casted,
	}
}

type InStringValidatorError struct {
	Required []string
	Actual   string
}

func (e InStringValidatorError) Error() string {
	return fmt.Sprintf("required one of [%s], actual: %s", strings.Join(e.Required, ","), e.Actual)
}
