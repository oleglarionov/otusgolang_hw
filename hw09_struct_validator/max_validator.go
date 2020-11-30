package hw09_struct_validator //nolint:golint,stylecheck,dupl

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type MaxValidator struct {
	maxValue int64
}

func (v *MaxValidator) Type() ValidatorType {
	return IntValidator
}

func (v *MaxValidator) TagName() string {
	return "max"
}

func (v *MaxValidator) Build(constraint string) error {
	maxValue, err := strconv.Atoi(constraint)
	if err != nil {
		return errors.Wrap(err, "conversion to int error")
	}

	v.maxValue = int64(maxValue)
	return nil
}

func (v *MaxValidator) Validate(value interface{}) error {
	casted, ok := value.(int64)
	if !ok {
		panic(fmt.Sprintf("can not cast %v to int64", value))
	}

	if casted > v.maxValue {
		return MaxValidatorError{
			RequiredMax: v.maxValue,
			ActualValue: casted,
		}
	}

	return nil
}

type MaxValidatorError struct {
	RequiredMax int64
	ActualValue int64
}

func (e MaxValidatorError) Error() string {
	return fmt.Sprintf("required max value: %d, actual value: %d", e.RequiredMax, e.ActualValue)
}
