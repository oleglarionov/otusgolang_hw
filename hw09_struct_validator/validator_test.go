package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "2a513df4-af0c-468e-9620-8a148d5eedf7",
				Name:   "Вася",
				Age:    30,
				Email:  "vpupkin@gmail.com",
				Role:   "admin",
				Phones: []string{"89998887766", "89995554433"},
				meta:   json.RawMessage(`{"field1": "value1"}`),
			},
			ValidationErrors(nil),
		},
		{
			User{
				ID:     "2a513df4-af0c-468e-9620-8a148d5eedf7",
				Name:   "Вася",
				Age:    30,
				Email:  "vpupkin@gmail.com",
				Role:   "admin",
				Phones: nil,
				meta:   json.RawMessage(`{"field1": "value1"}`),
			},
			ValidationErrors(nil),
		},
		{
			User{
				ID:     "2a513df4-af0c-468e-9620-8a148d5eedf7",
				Name:   "Вася",
				Age:    500,
				Email:  "vpupkin@gmail.com",
				Role:   "admin",
				Phones: nil,
				meta:   json.RawMessage(`{"field1": "value1"}`),
			},
			ValidationErrors{
				ValidationError{
					Key: "Age",
					Err: ValidationErrors{
						ValidationError{
							Key: "max",
							Err: MaxValidatorError{
								RequiredMax: 50,
								ActualValue: 500,
							},
						},
					},
				},
			},
		},
		{
			User{
				ID:     "1",
				Name:   "Вася",
				Age:    17,
				Email:  "v.pupkin@gmail.com",
				Role:   "manager",
				Phones: []string{"89998887766", "+79995554433"},
				meta:   json.RawMessage(`{"field1": "value1"}`),
			},
			ValidationErrors{
				ValidationError{
					Key: "ID",
					Err: ValidationErrors{
						ValidationError{
							Key: "len",
							Err: LengthValidatorError{
								Required: 36,
								Actual:   1,
							},
						},
					},
				},
				ValidationError{
					Key: "Age",
					Err: ValidationErrors{
						ValidationError{
							Key: "min",
							Err: MinValidatorError{
								RequiredMin: 18,
								ActualValue: 17,
							},
						},
					},
				},
				ValidationError{
					Key: "Email",
					Err: ValidationErrors{
						ValidationError{
							Key: "regexp",
							Err: RegexpValidatorError{
								RequiredRegexp: "^\\w+@\\w+\\.\\w+$",
								ActualValue:    "v.pupkin@gmail.com",
							},
						},
					},
				},
				ValidationError{
					Key: "Role",
					Err: ValidationErrors{
						ValidationError{
							Key: "in",
							Err: InStringValidatorError{
								Required: []string{"admin", "stuff"},
								Actual:   "manager",
							},
						},
					},
				},
				ValidationError{
					Key: "Phones",
					Err: ValidationErrors{
						ValidationError{
							Key: "1",
							Err: ValidationErrors{
								ValidationError{
									Key: "len",
									Err: LengthValidatorError{
										Required: 11,
										Actual:   12,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			App{
				Version: "1.0.0",
			},
			ValidationErrors(nil),
		},
		{
			App{
				Version: "10.0.0",
			},
			ValidationErrors{
				ValidationError{
					Key: "Version",
					Err: ValidationErrors{
						ValidationError{
							Key: "len",
							Err: LengthValidatorError{
								Required: 5,
								Actual:   6,
							},
						},
					},
				},
			},
		},
		{
			Token{
				Header:    nil,
				Payload:   []byte("Payload1"),
				Signature: []byte("Signature1"),
			},
			ValidationErrors(nil),
		},
		{
			Response{
				Code: 200,
				Body: "{}",
			},
			ValidationErrors(nil),
		},
		{
			Response{
				Code: 301,
				Body: "{}",
			},
			ValidationErrors{
				ValidationError{
					Key: "Code",
					Err: ValidationErrors{
						ValidationError{
							Key: "in",
							Err: InIntValidatorError{
								Required: []int64{200, 404, 500},
								Actual:   301,
							},
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
