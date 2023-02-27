package hw09structvalidator

import (
	"errors"
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
	ErrValueIsNotStruct            = errors.New("the value is not the structure")
	ErrCheckNotMatchType           = errors.New("the check does not match the type")
	ErrInWritingCheck              = errors.New("error in writing the check")
	ErrStringLengthLongerAllowed   = errors.New("Error validation: the string value length is longer than allowed")
	ErrValueLessMinimum            = errors.New("Error validation: the value is less than the minimum")
	ErrValueMoreMaximum            = errors.New("Error validation: the value is more than the maximum")
	ErrValueNotMatchPattern        = errors.New("Error validation: the string value does not match the pattern")
	ErrValueNotIncludedAllowedList = errors.New("Error validation: the value is not included in the allowed list")
)

func Validate(v interface{}) error {
	var vErrors ValidationErrors
	val := reflect.ValueOf(v)
	tVal := val.Type()
	if tVal.Kind() != reflect.Struct {
		return ErrValueIsNotStruct
	}
	for i := 0; i < tVal.NumField(); i++ {
		field := val.Field(i)
		if field.Type().Kind() != reflect.Slice {
			vErrors = rangeChecks(field, tVal.Field(i), vErrors)
		} else if field.Type().Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				vErrors = rangeChecks(elem, tVal.Field(i), vErrors)
			}
		}
	}
	return vErrors
}

func rangeChecks(field reflect.Value, tField reflect.StructField, vErrors ValidationErrors) ValidationErrors {
	tag := tField.Tag.Get("validate")
	if len(tag) == 0 {
		return vErrors
	}
	for _, ch := range strings.Split(tag, "|") {
		res := check(ch, field, tField.Name)
		if res != empty {
			vErrors = append(vErrors, res)
		}
	}
	return vErrors
}

func check(check string, field reflect.Value, name string) ValidationError {
	ch := strings.Split(check, ":")
	if len(ch) < 2 {
		return ValidationError{Field: name, Err: ErrInWritingCheck}
	}
	switch ch[0] {
	case "len":
		return checkLen(ch[1], field, name)
	case "regexp":
		return checkReg(ch[1], field, name)
	case "in":
		return checkIn(ch[1], field, name)
	case "min":
		return checkMin(ch[1], field, name)
	case "max":
		return checkMax(ch[1], field, name)
	default:
	}
	return empty
}

func checkLen(str string, field reflect.Value, name string) ValidationError {
	if field.Kind() != reflect.String {
		return ValidationError{Field: name, Err: ErrCheckNotMatchType}
	}
	l, err := strconv.Atoi(str)
	if err != nil {
		return ValidationError{Field: name, Err: ErrInWritingCheck}
	}
	if field.Len() > l {
		return ValidationError{Field: name, Err: ErrStringLengthLongerAllowed}
	}
	return empty
}

func checkIn(str string, field reflect.Value, name string) ValidationError {
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
				return ValidationError{Field: name, Err: ErrInWritingCheck}
			}
			if field.Int() == int64(i) {
				isIncluded = true
				break
			}
		}
	}
	if !isIncluded {
		return ValidationError{Field: name, Err: ErrValueNotIncludedAllowedList}
	}
	return empty
}

func checkMin(str string, field reflect.Value, name string) ValidationError {
	if field.Kind() != reflect.Int {
		return ValidationError{Field: name, Err: ErrCheckNotMatchType}
	}
	min, err := strconv.Atoi(str)
	if err != nil {
		return ValidationError{Field: name, Err: ErrInWritingCheck}
	}
	if int(field.Int()) < min {
		return ValidationError{Field: name, Err: ErrValueLessMinimum}
	}
	return empty
}

func checkMax(str string, field reflect.Value, name string) ValidationError {
	if field.Kind() != reflect.Int {
		return ValidationError{Field: name, Err: ErrCheckNotMatchType}
	}
	max, err := strconv.Atoi(str)
	if err != nil {
		return ValidationError{Field: name, Err: ErrInWritingCheck}
	}
	if int(field.Int()) > max {
		return ValidationError{Field: name, Err: ErrValueMoreMaximum}
	}
	return empty
}

func checkReg(str string, field reflect.Value, name string) ValidationError {
	if field.Kind() != reflect.String {
		return ValidationError{Field: name, Err: ErrCheckNotMatchType}
	}
	if mached, _ := regexp.MatchString(str, field.String()); !mached {
		return ValidationError{Field: name, Err: ErrValueNotMatchPattern}
	}
	return empty
}
