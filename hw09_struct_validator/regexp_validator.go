package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
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

func (v *RegexpValidator) Build(constraint string) error {
	r, err := regexp.Compile(constraint)
	if err != nil {
		return errors.Wrap(err, "regexp compile error")
	}

	v.r = r
	return nil
}

func (v *RegexpValidator) Validate(value interface{}) error {
	casted, ok := value.(string)
	if !ok {
		panic(fmt.Sprintf("can not cast %v to string", value))
	}

	if !v.r.MatchString(casted) {
		return RegexpValidatorError{
			RequiredRegexp: v.r.String(),
			ActualValue:    casted,
		}
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
