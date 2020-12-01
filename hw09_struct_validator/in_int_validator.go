package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type InIntValidator struct {
	values []int64
}

func (v *InIntValidator) Type() ValidatorType {
	return IntValidator
}

func (v *InIntValidator) TagName() string {
	return "in"
}

func (v *InIntValidator) Build(constraint string) error {
	strValues := strings.Split(constraint, ",")
	v.values = make([]int64, 0, len(strValues))
	for _, strValue := range strValues {
		intValue, err := strconv.Atoi(strValue)
		if err != nil {
			return errors.Wrap(err, "conversion to int error")
		}

		v.values = append(v.values, int64(intValue))
	}
	return nil
}

func (v *InIntValidator) Validate(value interface{}) error {
	casted, ok := value.(int64)
	if !ok {
		panic(fmt.Sprintf("can not cast %v to int64", value))
	}

	for _, value := range v.values {
		if casted == value {
			return nil
		}
	}

	return InIntValidatorError{
		Required: v.values,
		Actual:   casted,
	}
}

type InIntValidatorError struct {
	Required []int64
	Actual   int64
}

func (e InIntValidatorError) Error() string {
	strValues := make([]string, 0, len(e.Required))
	for _, intVal := range e.Required {
		strValues = append(strValues, strconv.FormatInt(intVal, 10))
	}

	return fmt.Sprintf("required one of [%s], actual: %d", strings.Join(strValues, ","), e.Actual)
}
