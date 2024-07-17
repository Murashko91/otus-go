package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/slices"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errors := make([]string, 0, len(v))

	for _, err := range v {
		errors = append(errors, fmt.Sprintf("field: %s, message: %s", err.Field, err.Err))
	}
	return strings.Join(errors, "\n")
}

func Validate(v interface{}) error {
	rValueStruct := reflect.ValueOf(v)
	if rValueStruct.Type().Kind() != reflect.Struct {
		return nil
	}
	rTypeStruct := rValueStruct.Type()
	numbFields := rTypeStruct.NumField()
	result := make(ValidationErrors, 0)

	for i := 0; i < numbFields; i++ {
		rField := rTypeStruct.Field(i)
		rFieldValue := rValueStruct.Field(i)

		rTag := rField.Tag
		rFieldType := rField.Type
		rKind := rFieldType.Kind().String()

		validateTag, ok := rTag.Lookup("validate")
		if ok {
			valErrors, err := processKind(rKind, validateTag, rFieldValue, rField.Name)
			if err != nil {
				return err
			}

			result = append(result, valErrors...)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return fmt.Errorf("validation errors: [%w]", result)
}

func processKind(
	rKind string,
	validateTag string,
	rFieldValue reflect.Value,
	rFieldName string,
) (ValidationErrors, error) {
	result := make(ValidationErrors, 0)

	switch rKind {
	case "int":
		rr, err := validateInt(validateTag, rFieldValue, rFieldName)
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			result = append(result, rr...)
		}

	case "string":
		rr, err := validateString(validateTag, rFieldValue, rFieldName)
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			result = append(result, rr...)
		}
	case "slice":
		sliceKind := rFieldValue.Type().Elem().Kind().String()

		if slices.Contains([]string{reflect.String.String(), reflect.Int.String(), reflect.Slice.String()}, sliceKind) {
			for i := 0; i < rFieldValue.Len(); i++ {
				valErrors, err := processKind(sliceKind, validateTag, rFieldValue.Index(i), rFieldName)
				if err != nil {
					return nil, err
				}
				result = append(result, valErrors...)
			}
		}
	}
	return result, nil
}

func validateInt(vTag string, refValue reflect.Value, fieldName string) (ValidationErrors, error) {
	intValue := int(refValue.Int())

	result := make(ValidationErrors, 0)

	tags := strings.Split(vTag, "|")
	for _, tag := range tags {
		switch strings.Split(tag, ":")[0] {
		case "min":
			validationVal, err := strconv.Atoi(tag[4:])
			if err != nil {
				return nil, err
			}
			if intValue < validationVal {
				result = append(result, generateValidationError(tag, intValue, fieldName))
			}
		case "max":
			validationVal, err := strconv.Atoi(tag[4:])
			if err != nil {
				return nil, err
			}
			if intValue > validationVal {
				result = append(result, generateValidationError(tag, intValue, fieldName))
			}

		case "in":
			strVal := strconv.Itoa(intValue)
			strValues := strings.Split(tag[3:], ",")
			if !slices.Contains(strValues, strVal) {
				result = append(result, generateValidationError(tag, intValue, fieldName))
			}
		}
	}

	return result, nil
}

func validateString(vTag string, refValue reflect.Value, fieldName string) (ValidationErrors, error) {
	stringValue := refValue.String()

	result := make(ValidationErrors, 0)

	tags := strings.Split(vTag, "|")
	for _, tag := range tags {
		switch strings.Split(tag, ":")[0] {
		case "len":
			validationVal, err := strconv.Atoi(tag[4:])
			if err != nil {
				return nil, err
			}

			if utf8.RuneCountInString(stringValue) != validationVal {
				result = append(result, generateValidationError(tag, stringValue, fieldName))
			}

		case "regexp":
			validationVal := tag[7:]

			valRegexp, err := regexp.Compile(validationVal)
			if err != nil {
				return nil, err
			}

			if !valRegexp.MatchString(stringValue) {
				result = append(result, generateValidationError(tag, stringValue, fieldName))
			}

		case "in":

			strValues := strings.Split(tag[3:], ",")
			if !slices.Contains(strValues, stringValue) {
				result = append(result, generateValidationError(tag, stringValue, fieldName))
			}
		default:
			// Ignore validation if not implemented?
		}
	}

	return result, nil
}

func generateValidationError(tag string, value interface{}, fieldName string) ValidationError {
	return ValidationError{
		Field: fieldName,
		Err:   fmt.Errorf("tag: %s, value: %v", tag, value),
	}
}
