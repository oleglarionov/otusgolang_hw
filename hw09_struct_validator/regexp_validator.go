package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"regexp"
)

type RegexpValidator struct {
	r *regexp.Regexp
}

func (v *RegexpValidator) Type() ValidatorType {
	return StringValidator
}

func (v *RegexpValidator) TagName() string {
	return "regexp"
}

func (v *RegexpValidator) Build(constraint string) {
	r, err := regexp.Compile(constraint)
	if err != nil {
		panic(err)
	}
	v.r = r
}

func (v *RegexpValidator) Validate(value interface{}) error {
	switch casted := value.(type) {
	case string:
		if !v.r.MatchString(casted) {
			return RegexpValidatorError{
				RequiredRegexp: v.r.String(),
				ActualValue:    casted,
			}
		}
	default:
		panic("invalid value type")
	}

	return nil
}

type RegexpValidatorError struct {
	RequiredRegexp string
	ActualValue    string
}

func (e RegexpValidatorError) Error() string {
	return fmt.Sprintf("required regexp: %s, actual value: %s", e.RequiredRegexp, e.ActualValue)
}
