package hw09_struct_validator //nolint:golint,stylecheck
import (
	"encoding/json"
	"strings"
)

type ValidationError struct {
	Key string
	Err error
}

func (v ValidationError) MarshalJSON() ([]byte, error) {
	marshal, _ := json.Marshal(map[string]string{
		"Key": v.Key,
		"Err": v.Err.Error(),
	})
	return marshal, nil
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	keyStack := make([]string, 0)
	errsMap := make(map[string]string)
	v.fillErrorsMap(keyStack, errsMap)

	marshal, _ := json.Marshal(errsMap)

	return string(marshal)
}

func (v *ValidationErrors) fillErrorsMap(keyStack []string, errsMap map[string]string) {
	for _, vErr := range *v {
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

func Validate(v interface{}) error {
	validator := MakeStructValidator([]Validator{
		&LengthValidator{},
		&RegexpValidator{},
		&InStringValidator{},
		&MinValidator{},
		&MaxValidator{},
		&InIntValidator{},
	})

	return validator.Validate(v)
}
