package hw09_struct_validator //nolint:golint,stylecheck,dupl

import (
	"fmt"
	"strconv"
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

func (v *MaxValidator) Build(constraint string) {
	maxValue, err := strconv.Atoi(constraint)
	if err != nil {
		panic(err)
	}

	v.maxValue = int64(maxValue)
}

func (v *MaxValidator) Validate(value interface{}) error {
	casted := value.(int64)
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
