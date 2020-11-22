package hw09_struct_validator //nolint:golint,stylecheck,dupl

import (
	"fmt"
	"strconv"
)

type MinValidator struct {
	minValue int64
}

func (v *MinValidator) Type() ValidatorType {
	return IntValidator
}

func (v *MinValidator) TagName() string {
	return "min"
}

func (v *MinValidator) Build(constraint string) {
	minValue, err := strconv.Atoi(constraint)
	if err != nil {
		panic(err)
	}

	v.minValue = int64(minValue)
}

func (v *MinValidator) Validate(value interface{}) error {
	casted := value.(int64)
	if casted < v.minValue {
		return MinValidatorError{
			RequiredMin: v.minValue,
			ActualValue: casted,
		}
	}

	return nil
}

type MinValidatorError struct {
	RequiredMin int64
	ActualValue int64
}

func (e MinValidatorError) Error() string {
	return fmt.Sprintf("required min value: %d, actual value: %d", e.RequiredMin, e.ActualValue)
}
