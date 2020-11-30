package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type LengthValidator struct {
	lenValue int
}

func (v *LengthValidator) Build(constraint string) error {
	value, err := strconv.Atoi(constraint)
	if err != nil {
		return errors.Wrap(err, "conversion to int error")
	}

	v.lenValue = value
	return nil
}

func (v *LengthValidator) Validate(value interface{}) error {
	casted, ok := value.(string)
	if !ok {
		panic(fmt.Sprintf("can not cast %v to string", value))
	}

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
