package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ValidatorType int

const (
	StringValidator ValidatorType = iota
	IntValidator
	ValidationTag = "validate"
)

type StructValidator struct {
	validators map[ValidatorType]map[string]Validator
}

type Validator interface {
	Type() ValidatorType
	TagName() string
	Build(constraint string) error
	Validate(value interface{}) error
}

func MakeStructValidator(validators []Validator) (*StructValidator, error) {
	vMap := make(map[ValidatorType]map[string]Validator)
	for _, v := range validators {
		vType := v.Type()
		vTagName := v.TagName()

		vByType, ok := vMap[vType]
		if !ok {
			vMap[vType] = make(map[string]Validator)
			vByType = vMap[vType]
		}

		_, ok = vByType[vTagName]
		if ok {
			return nil, fmt.Errorf("validator with TagName %s duplicated", vTagName)
		}

		vByType[vTagName] = v
	}

	return &StructValidator{validators: vMap}, nil
}

func (v *StructValidator) Validate(value interface{}) (ValidationErrors, error) {
	rv := reflect.ValueOf(value)
	t := rv.Type()

	if t.Kind() != reflect.Struct {
		return nil, errors.New("input value must be a structure")
	}

	validationErrs := ValidationErrors{}
	for i := 0; i < rv.NumField(); i++ {
		vf := rv.Field(i)
		if !vf.CanInterface() {
			continue
		}

		tf := t.Field(i)
		validateTag, ok := tf.Tag.Lookup(ValidationTag)
		if !ok {
			continue
		}

		var validators []Validator

		elKind := getElKind(tf.Type)

		var err error
		switch elKind { //nolint:exhaustive
		case reflect.Int:
			validators, err = v.buildValidators(IntValidator, validateTag)
		case reflect.String:
			validators, err = v.buildValidators(StringValidator, validateTag)
		default:
			return nil, errors.New("unexpected type: " + elKind.String())
		}

		if err != nil {
			return nil, err
		}

		fieldErrs, err := validateFieldWithValidators(validators, vf)
		if err != nil {
			return nil, err
		}

		if len(fieldErrs) > 0 {
			validationErrs = append(validationErrs, ValidationError{
				Key: tf.Name,
				Err: fieldErrs,
			})
		}
	}

	if len(validationErrs) == 0 {
		return nil, nil
	}

	return validationErrs, nil
}

func (v *StructValidator) buildValidators(vType ValidatorType, validateTag string) ([]Validator, error) {
	validatorPairs := strings.Split(validateTag, "|")
	fieldValidators := make([]Validator, 0, len(validatorPairs))

	for _, validatorPair := range validatorPairs {
		validatorParts := strings.SplitN(validatorPair, ":", 2)
		if len(validatorParts) != 2 {
			return nil, errors.New("invalid validate tag: " + validatorPair)
		}

		validatorTagName := validatorParts[0]
		constraint := validatorParts[1]

		validator, ok := v.validators[vType][validatorTagName]
		if !ok {
			return nil, errors.New("unknown validator with tagName: " + validatorTagName)
		}

		err := validator.Build(constraint)
		if err != nil {
			return nil, errors.Wrap(err, "build validator error")
		}

		fieldValidators = append(fieldValidators, validator)
	}

	return fieldValidators, nil
}

func getElKind(t reflect.Type) reflect.Kind {
	kind := t.Kind()
	switch kind { //nolint:exhaustive
	case reflect.Slice:
		return getElKind(t.Elem())
	default:
		return kind
	}
}

func validateFieldWithValidators(validators []Validator, value reflect.Value) (ValidationErrors, error) {
	var fieldErrs ValidationErrors

	kind := value.Kind()
	switch kind { //nolint:exhaustive
	case reflect.Int:
		fieldErrs = validate(validators, value.Int())
	case reflect.String:
		fieldErrs = validate(validators, value.String())
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			elErrors, err := validateFieldWithValidators(validators, value.Index(i))
			if err != nil {
				return nil, err
			}

			if len(elErrors) > 0 {
				fieldErrs = append(fieldErrs, ValidationError{
					Key: strconv.Itoa(i),
					Err: elErrors,
				})
			}
		}
	default:
		return nil, errors.New("unexpected type: " + kind.String())
	}

	return fieldErrs, nil
}

func validate(validators []Validator, value interface{}) ValidationErrors {
	errs := make(ValidationErrors, 0, len(validators))
	for _, validator := range validators {
		err := validator.Validate(value)
		if err != nil {
			errs = append(errs, ValidationError{
				Key: validator.TagName(),
				Err: err,
			})
		}
	}

	return errs
}
