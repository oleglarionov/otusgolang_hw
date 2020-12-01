package hw09_struct_validator //nolint:golint,stylecheck,dupl

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
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

func (v *MinValidator) Build(constraint string) error {
	minValue, err := strconv.Atoi(constraint)
	if err != nil {
		return errors.Wrap(err, "conversion to int error")
	}

	v.minValue = int64(minValue)
	return nil
}

func (v *MinValidator) Validate(value interface{}) error {
	casted, ok := value.(int64)
	if !ok {
		panic(fmt.Sprintf("can not cast %v to int64", value))
	}

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
