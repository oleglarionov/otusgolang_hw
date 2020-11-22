package hw09_struct_validator //nolint:golint,stylecheck

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

type ValidatorType int

const (
	StringValidator ValidatorType = iota
	IntValidator
)

type StructValidator struct {
	validators map[ValidatorType]map[string]Validator
}

type Validator interface {
	Type() ValidatorType
	TagName() string
	Build(constraint string)
	Validate(value interface{}) error
}

func MakeStructValidator(validators []Validator) *StructValidator {
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
			log.Fatalf("Validator with TagName %s duplicated", vTagName)
		}

		vByType[vTagName] = v
	}

	return &StructValidator{validators: vMap}
}

func (v *StructValidator) Validate(value interface{}) ValidationErrors {
	rv := reflect.ValueOf(value)
	t := rv.Type()

	errs := make(ValidationErrors, 0)
	for i := 0; i < rv.NumField(); i++ {
		vf := rv.Field(i)
		if !vf.CanInterface() {
			continue
		}

		tf := t.Field(i)
		validateTag, ok := tf.Tag.Lookup("validate")
		if !ok {
			continue
		}

		var validators []Validator

		elKind := getElKind(tf.Type)
		switch elKind { //nolint:exhaustive
		case reflect.Int:
			validators = v.buildValidators(IntValidator, validateTag)
		case reflect.String:
			validators = v.buildValidators(StringValidator, validateTag)
		default:
			log.Fatal("unexpected type")
		}

		fieldErrs := validateFieldWithValidators(validators, vf)

		if len(fieldErrs) > 0 {
			errs = append(errs, ValidationError{
				Key: tf.Name,
				Err: fieldErrs,
			})
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func (v *StructValidator) buildValidators(vType ValidatorType, validateTag string) []Validator {
	validatorPairs := strings.Split(validateTag, "|")
	fieldValidators := make([]Validator, 0, len(validatorPairs))

	for _, validatorPair := range validatorPairs {
		validatorParts := strings.SplitN(validatorPair, ":", 2)
		if len(validatorParts) != 2 {
			log.Fatal("invalid validate tag")
		}

		validatorTagName := validatorParts[0]
		constraint := validatorParts[1]

		validator, ok := v.validators[vType][validatorTagName]
		if !ok {
			// log.Fatalf("unknown validator with tagName %s", validatorTagName)
			continue
		}

		validator.Build(constraint)

		fieldValidators = append(fieldValidators, validator)
	}

	return fieldValidators
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

func validateFieldWithValidators(validators []Validator, value reflect.Value) ValidationErrors {
	fieldErrs := make(ValidationErrors, 0)

	switch value.Kind() { //nolint:exhaustive
	case reflect.Int:
		fieldErrs = validate(validators, value.Int())
	case reflect.String:
		fieldErrs = validate(validators, value.String())
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			elErrors := validateFieldWithValidators(validators, value.Index(i))
			if len(elErrors) > 0 {
				fieldErrs = append(fieldErrs, ValidationError{
					Key: strconv.Itoa(i),
					Err: elErrors,
				})
			}
		}
	default:
		log.Fatal("unexpected type")
	}

	return fieldErrs
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
