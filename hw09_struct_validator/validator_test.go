package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp/syntax"
	"strconv"
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

func getInvalidUserError() ValidationErrors {
	result := make(ValidationErrors, 0, 6)

	result = append(result, generateValidationError("len:36", 1, "ID"))
	result = append(result, generateValidationError("min:18", 1, "Age"))
	result = append(result, generateValidationError("regexp:^\\w+@\\w+\\.\\w+$", "invalid", "Email"))
	result = append(result, generateValidationError("in:admin,stuff", "balbes", "Role"))
	result = append(result, generateValidationError("len:11", "invalid2", "Phones"))

	return result
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "550e8400-e29b-41d4-a716-446655440000",
				Name:   "Boddah",
				Age:    27,
				Email:  "hello@test.com",
				Role:   "admin",
				Phones: []string{"11111111111", "+37529-test"},
				meta:   []byte("{}"),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "1",
				Age:    1,
				Email:  "invalid",
				Role:   "balbes",
				Phones: []string{"11111111111", "invalid2"},
			},
			expectedErr: getInvalidUserError(),
		},
		{
			in:          "not a struct",
			expectedErr: nil,
		},
		{
			in: struct {
				Version string `validate:"regexp:[\\]{hello}"`
			}{},
			expectedErr: &syntax.Error{
				Code: "missing closing ]",
				Expr: "[\\]{hello}",
			},
		},

		{
			in: struct {
				Count int `validate:"min:error"`
			}{},
			expectedErr: &strconv.NumError{
				Func: "Atoi",
				Num:  "error",
				Err:  errors.New("invalid syntax"),
			},
		},
		{
			in: struct {
				Count int `validate:"max:error"`
			}{},
			expectedErr: &strconv.NumError{
				Func: "Atoi",
				Num:  "error",
				Err:  errors.New("invalid syntax"),
			},
		},
		{
			in: struct {
				Count string `validate:"len:error"`
			}{},
			expectedErr: &strconv.NumError{
				Func: "Atoi",
				Num:  "error",
				Err:  errors.New("invalid syntax"),
			},
		},
		{
			in: struct {
				Count string `validate:"notImplemented:error"`
			}{},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			// i := i
			t.Parallel()

			validateResult := Validate(tt.in)
			// Place your code here.
			fmt.Println(validateResult)
			if validateResult == nil {
				require.Equal(t, validateResult, tt.expectedErr)
			} else {
				var valErr ValidationErrors
				if errors.As(validateResult, &valErr) {
					require.Equal(t, errors.Unwrap(validateResult), tt.expectedErr)
				} else {
					require.Equal(t, validateResult, tt.expectedErr)
				}
			}
		})
	}
}
