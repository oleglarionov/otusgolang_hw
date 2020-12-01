package hw09_struct_validator //nolint:golint,stylecheck
import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type ValidationError struct {
	Key string
	Err error
}

func (v ValidationError) MarshalJSON() ([]byte, error) {
	marshal, err := json.Marshal(map[string]string{
		"Key": v.Key,
		"Err": v.Err.Error(),
	})
	return marshal, errors.Wrap(err, "marshaling error")
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var keyStack []string
	errsMap := make(map[string]string)
	v.fillErrorsMap(keyStack, errsMap)

	marshal, err := json.Marshal(errsMap)
	if err != nil {
		panic(err)
	}

	return string(marshal)
}

func (v ValidationErrors) fillErrorsMap(keyStack []string, errsMap map[string]string) {
	for _, vErr := range v {
		keyStack = append(keyStack, vErr.Key)

		nestedV, ok := vErr.Err.(ValidationErrors) //nolint:errorlint
		if ok {
			nestedV.fillErrorsMap(keyStack, errsMap)
		} else {
			errKey := strings.Join(keyStack, ".")
			errsMap[errKey] = vErr.Err.Error()
		}

		keyStack = keyStack[0 : len(keyStack)-1]
	}
}

var validator *StructValidator

func init() {
	var err error
	validator, err = MakeStructValidator([]Validator{
		&LengthValidator{},
		&RegexpValidator{},
		&InStringValidator{},
		&MinValidator{},
		&MaxValidator{},
		&InIntValidator{},
	})

	if err != nil {
		panic(err)
	}
}

func Validate(v interface{}) (ValidationErrors, error) {
	validationErrors, err := validator.Validate(v)
	if err != nil {
		return nil, err
	}

	return validationErrors, nil
}
