package main

import (
	"encoding/json"
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

func (v ValidationErrors) Error() string {
	panic("implement me")
}

var empty = ValidationError{}

var (
	ErrCheckNotMatchType           = errors.New("the check does not match the type")
	ErrInWritingCheck              = errors.New("error in writing the check")
	ErrStringLengthLongerAllowed   = errors.New("Error validation: the string value length is longer than allowed")
	ErrValueLessMinimum            = errors.New("Error validation: the value is less than the minimum")
	ErrValueMoreMaximum            = errors.New("Error validation: the value is more than the maximum")
	ErrValueNotMatchPattern        = errors.New("Error validation: the string value does not match the pattern")
	ErrValueNotIncludedAllowedList = errors.New("Error validation: the value is not included in the allowed list")
)

type ValidationErrors []ValidationError

type UserRole string

type User struct {
	ID     string `json:"id" validate:"len:32"`
	Code   int    `validate:"in:25,36,87"`
	Name   string
	Age    int      `validate:"min:18|max:50"`
	Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	Role   UserRole `validate:"in:admin,stuff"`
	Phones []string `validate:"len:11"`
	meta   json.RawMessage
}

func main() {
	user := User{
		"12234567876543212345",
		36,
		"name",
		25,
		"1@ya.ru",
		UserRole("admin"),
		[]string{"12312131523"},
		[]byte("1dfdfsd"),
	}
	err := Validate(user)
	if err != nil {
		panic("Panic!")
	}
}

func Validate(v interface{}) error {
	var vErrors ValidationErrors
	val := reflect.ValueOf(v)
	tVal := val.Type()
	for i := 0; i < tVal.NumField(); i++ {
		field := val.Field(i)
		if field.Type().Kind() != reflect.Slice {
			rangeChecks(field, tVal.Field(i), vErrors)
		} else if field.Type().Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				rangeChecks(elem, tVal.Field(i), vErrors)
			}
		}
	}
	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func rangeChecks(field reflect.Value, tField reflect.StructField, vErrors ValidationErrors) {
	tag := tField.Tag.Get("validate")
	if len(tag) == 0 {
		return
	}
	for _, ch := range strings.Split(tag, "|") {
		res := check(ch, field, tField)
		if res != empty {
			vErrors = append(vErrors, res)
			fmt.Printf("%s: %s\n", res.Field, res.Err.Error())
		}
	}
}

func check(check string, field reflect.Value, tField reflect.StructField) ValidationError {
	name := tField.Name
	ch := strings.Split(check, ":")
	if len(ch) < 2 {
		return ValidationError{Field: name, Err: ErrInWritingCheck}
	}
	switch ch[0] {
	case "len":
		if field.Kind() != reflect.String {
			return ValidationError{Field: name, Err: ErrCheckNotMatchType}
		}
		l, err := strconv.Atoi(ch[1])
		if err != nil {
			return ValidationError{Field: name, Err: ErrInWritingCheck}
		}
		if field.Len() > l {
			return ValidationError{Field: name, Err: ErrStringLengthLongerAllowed}
		}
	case "regexp":
		if field.Kind() != reflect.String {
			return ValidationError{Field: name, Err: ErrCheckNotMatchType}
		}
		if mached, _ := regexp.MatchString(ch[1], field.String()); !mached {
			return ValidationError{Field: name, Err: ErrValueNotMatchPattern}
		}
	case "in":
		values := strings.Split(ch[1], ",")
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
	case "min":
		if field.Kind() != reflect.Int {
			return ValidationError{Field: name, Err: ErrCheckNotMatchType}
		}
		min, err := strconv.Atoi(ch[1])
		if err != nil {
			return ValidationError{Field: name, Err: ErrInWritingCheck}
		}
		if int(field.Int()) < min {
			return ValidationError{Field: name, Err: ErrValueLessMinimum}
		}
	case "max":
		if field.Kind() != reflect.Int {
			return ValidationError{Field: name, Err: ErrCheckNotMatchType}
		}
		max, err := strconv.Atoi(ch[1])
		if err != nil {
			return ValidationError{Field: name, Err: ErrInWritingCheck}
		}
		if int(field.Int()) > max {
			return ValidationError{Field: name, Err: ErrValueMoreMaximum}
		}
	default:
	}
	return empty
}
