package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var empty = ValidationError{}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

var (
	ErrStringLengthLongerAllowed   = errors.New("Error validation: the string value length is longer than allowed")
	ErrValueLessMinimum            = errors.New("Error validation: the value is less than the minimum")
	ErrValueMoreMaximum            = errors.New("Error validation: the value is more than the maximum")
	ErrValueNotMatchPattern        = errors.New("Error validation: the string value does not match the pattern")
	ErrValueNotIncludedAllowedList = errors.New("Error validation: the value is not included in the allowed list")
)

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	tVal := val.Type()
	if tVal.Kind() != reflect.Struct {
		return fmt.Errorf("the value is not the structure")
	}
	var vErrors ValidationErrors
	vErrors, err := validate(val, tVal, vErrors)
	if err != nil {
		return err
	}
	return vErrors
}

func validate(val reflect.Value, tVal reflect.Type, vErrors ValidationErrors) (ValidationErrors, error) {
	var err error
	for i := 0; i < tVal.NumField(); i++ {
		field := val.Field(i)
		if field.Type().Kind() == reflect.Int ||
			field.Type().Kind() == reflect.String ||
			field.Type().Kind() == reflect.Struct {
			vErrors, err = rangeChecks(field, tVal.Field(i), vErrors)
			if err != nil {
				return nil, err
			}
		} else if field.Type().Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				vErrors, err = rangeChecks(elem, tVal.Field(i), vErrors)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return vErrors, nil
}

func rangeChecks(field reflect.Value, tField reflect.StructField, vErrors ValidationErrors) (ValidationErrors, error) {
	tag := tField.Tag.Get("validate")
	if len(tag) == 0 {
		return vErrors, nil
	}
	for _, ch := range strings.Split(tag, "|") {
		var err error
		vErrors, err = check(ch, field, tField.Name, vErrors)
		if err != nil {
			return vErrors, err
		}
	}
	return vErrors, nil
}

func check(check string, field reflect.Value, name string, vErrors ValidationErrors) (ValidationErrors, error) {
	ch := strings.Split(check, ":")
	if len(ch) < 2 || len(ch) == 1 && ch[0] != "nested" {
		return vErrors, fmt.Errorf("validation string is not correct")
	}
	var validationErr ValidationError
	var err error
	switch ch[0] {
	case "len":
		validationErr, err = checkLen(ch[1], field, name)
	case "regexp":
		validationErr, err = checkReg(ch[1], field, name)
	case "in":
		validationErr, err = checkIn(ch[1], field, name)
	case "min":
		validationErr, err = checkMin(ch[1], field, name)
	case "max":
		validationErr, err = checkMax(ch[1], field, name)
	case "nested":
		validate(field, field.Type(), vErrors)
	default:
	}
	if err != nil {
		return nil, err
	}
	if validationErr != empty {
		vErrors = append(vErrors, validationErr)
	}
	return vErrors, nil
}

func checkLen(str string, field reflect.Value, name string) (ValidationError, error) {
	if field.Kind() != reflect.String {
		return empty, fmt.Errorf("field %s: the check does not match the type", name)
	}
	l, err := strconv.Atoi(str)
	if err != nil {
		return empty, fmt.Errorf("field %s: error is %w", name, err)
	}
	if field.Len() > l {
		return ValidationError{Field: name, Err: ErrStringLengthLongerAllowed}, nil
	}
	return empty, nil
}

func checkIn(str string, field reflect.Value, name string) (ValidationError, error) {
	values := strings.Split(str, ",")
	isIncluded := false
	for _, val := range values {
		if field.Kind() == reflect.String {
			if strings.Compare(field.String(), val) == 0 {
				isIncluded = true
				break
			}
		} else if field.Kind() == reflect.Int {
			i, err := strconv.Atoi(val)
			if err != nil {
				return empty, fmt.Errorf("field %s: error is %w", name, err)
			}
			if field.Int() == int64(i) {
				isIncluded = true
				break
			}
		}
	}
	if !isIncluded {
		return ValidationError{Field: name, Err: ErrValueNotIncludedAllowedList}, nil
	}
	return empty, nil
}

func checkMin(str string, field reflect.Value, name string) (ValidationError, error) {
	if field.Kind() != reflect.Int {
		return empty, fmt.Errorf("field %s: the check does not match the type", name)
	}
	min, err := strconv.Atoi(str)
	if err != nil {
		return empty, fmt.Errorf("field %s: error is %w", name, err)
	}
	if int(field.Int()) < min {
		return ValidationError{Field: name, Err: ErrValueLessMinimum}, nil
	}
	return empty, nil
}

func checkMax(str string, field reflect.Value, name string) (ValidationError, error) {
	if field.Kind() != reflect.Int {
		return empty, fmt.Errorf("field %s: the check does not match the type", name)
	}
	max, err := strconv.Atoi(str)
	if err != nil {
		return empty, fmt.Errorf("field %s: error is %w", name, err)
	}
	if int(field.Int()) > max {
		return ValidationError{Field: name, Err: ErrValueMoreMaximum}, nil
	}
	return empty, nil
}

func checkReg(str string, field reflect.Value, name string) (ValidationError, error) {
	if field.Kind() != reflect.String {
		return empty, fmt.Errorf("field %s: the check does not match the type", name)
	}
	if mached, _ := regexp.MatchString(str, field.String()); !mached {
		return ValidationError{Field: name, Err: ErrValueNotMatchPattern}, nil
	}
	return empty, nil
}
