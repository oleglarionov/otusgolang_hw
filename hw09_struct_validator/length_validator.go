package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"strconv"
)

type LengthValidator struct {
	lenValue int
}

func (v *LengthValidator) Build(constraint string) {
	value, err := strconv.Atoi(constraint)
	if err != nil {
		panic(err)
	}

	v.lenValue = value
}

func (v *LengthValidator) Validate(value interface{}) error {
	casted := value.(string)

	actualLen := len(casted)
	if actualLen != v.lenValue {
		return LengthValidatorError{
			Required: v.lenValue,
			Actual:   actualLen,
		}
	}

	return nil
}

func (v *LengthValidator) Type() ValidatorType {
	return StringValidator
}

func (v *LengthValidator) TagName() string {
	return "len"
}

type LengthValidatorError struct {
	Required int
	Actual   int
}

func (l LengthValidatorError) Error() string {
	return fmt.Sprintf("required length: %d, actual length: %d", l.Required, l.Actual)
}
