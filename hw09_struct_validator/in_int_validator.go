package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"strconv"
	"strings"
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

func (v *InIntValidator) Build(constraint string) {
	strValues := strings.Split(constraint, ",")
	v.values = make([]int64, 0, len(strValues))
	for _, strValue := range strValues {
		intValue, err := strconv.Atoi(strValue)
		if err != nil {
			panic(err)
		}

		v.values = append(v.values, int64(intValue))
	}
}

func (v *InIntValidator) Validate(value interface{}) error {
	casted := value.(int64)
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
