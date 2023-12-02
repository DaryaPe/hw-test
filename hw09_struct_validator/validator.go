package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validateTag        = "validate"
	validateTagDelim   = "|"
	validateCheckDelim = ":"
	validateListDelim  = ","
	nestedTag          = "nested"
)

const (
	checkLen    = "len"
	checkRegexp = "regexp"
	checkIn     = "in"
	checkMin    = "min"
	checkMax    = "max"
)

var (
	msgExpectedStruct = "expected struct but returned %s value"
	msgExpectedValue  = "%w: tag '%v' have not value"
	msgUndefinedTag   = "%w: tag '%v' is undefined"
	msgIncorrectTag   = "%w: incorrect value type for tag '%v' value '%v'"
	msgUnexpectedType = "%w: unexpected type for validation '%v' function"

	msgWrongLen    = "the length of the value must be '%d' not a '%d'"
	msgWrongRegexp = "value does not math pattern '%v'"
	msgLessMin     = "value '%d' is less than the limit, must be greater than '%d'"
	msgGreaterMax  = "value '%d' is greater than the limit, must be less than '%d'"
	msgWrongIn     = "value '%v' does not math any values from list '%v'"

	errInternal   = fmt.Errorf("internal error")
	errSystem     = fmt.Errorf("system error")
	errValidation = fmt.Errorf("validation error")
)

// ValidationError ошибки валидации.
type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result string
	for i := range v {
		result += fmt.Sprintf("%v: %v\n", v[i].Field, v[i].Err)
	}
	return result
}

// Validator объект для валидации структуры.
type Validator struct {
	validationErrs ValidationErrors // Ошибки валидации
}

// Struct валидирует переданную структуру.
func (v *Validator) Struct(i interface{}) error {
	if i == nil {
		return fmt.Errorf(msgExpectedStruct, "nil")
	}

	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Struct {
		return fmt.Errorf(msgExpectedStruct, value.Kind())
	}
	return v.validateStruct(value, value.Type())
}

// validateStruct валидирует структуру.
func (v *Validator) validateStruct(value reflect.Value, valueType reflect.Type) error {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldStruct := valueType.Field(i)

		tags, ok := fieldStruct.Tag.Lookup(validateTag)
		if !ok {
			continue
		}

		if err := v.Field(fieldStruct.Name, field, tags); err != nil {
			return err
		}
	}
	return nil
}

// Field валидирует переданное поле.
func (v *Validator) Field(name string, field reflect.Value, tags string) error {
	var err error

	switch {
	case field.Kind() == reflect.Struct && isNested(tags):
		if err = v.validateStruct(field, field.Type()); err != nil {
			return err
		}
	case field.Kind() == reflect.Slice && isNested(tags):
		for j := 0; j < field.Len(); j++ {
			if err = v.validateStruct(field.Index(j), field.Index(j).Type()); err != nil {
				return err
			}
		}
	case field.Kind() != reflect.Struct && field.Kind() != reflect.Slice:
		if err = v.validateField(name, field, tags); err != nil {
			return err
		}
	case field.Kind() != reflect.Struct && field.Kind() == reflect.Slice:
		for j := 0; j < field.Len(); j++ {
			if err = v.validateField(name, field.Index(j), tags); err != nil {
				return err
			}
		}
	}
	return nil
}

// isNested возвращает признак, что вложенная структура требует валидации.
func isNested(tags string) bool {
	return tags == nestedTag
}

// validateField валидирует атомарное поле.
func (v *Validator) validateField(name string, value reflect.Value, tags string) error {
	checks := strings.Split(tags, validateTagDelim)
	for _, check := range checks {
		var err error
		metadata := strings.Split(check, validateCheckDelim)
		if len(metadata) != 2 {
			return fmt.Errorf(msgExpectedValue, errInternal, metadata[0])
		}

		nameCheck, valueCheck := metadata[0], metadata[1]
		switch nameCheck {
		case checkLen:
			err = validateEqualLen(value, valueCheck)
		case checkRegexp:
			err = validateRegexp(value, valueCheck)
		case checkIn:
			err = validateIn(value, valueCheck)
		case checkMin:
			err = validateMin(value, valueCheck)
		case checkMax:
			err = validateMax(value, valueCheck)
		default:
			err = fmt.Errorf(msgUndefinedTag, errInternal, metadata[0])
		}
		if err != nil && errors.Is(err, errInternal) {
			return err
		}
		if err != nil {
			v.validationErrs = append(v.validationErrs, ValidationError{
				Field: name,
				Err:   err,
			})
		}
	}
	return nil
}

// Errors возвращает ошибки, обнаруженные в ходе валидации.
func (v *Validator) Errors() error {
	if len(v.validationErrs) > 0 {
		return fmt.Errorf("%w:\n%v", errValidation, v.validationErrs.Error())
	}
	return nil
}

// validateEqualLen проверяет равенство длины.
func validateEqualLen(value reflect.Value, check string) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		long, err := strconv.Atoi(check)
		if err != nil {
			return fmt.Errorf(msgIncorrectTag, errInternal, checkLen, value)
		}
		if value.Len() != long {
			return fmt.Errorf(msgWrongLen, long, value.Len())
		}
	default:
		return fmt.Errorf(msgUnexpectedType, errInternal, checkLen)
	}

	return nil
}

// validateRegexp проверяет на соответствие регулярному выражению.
func validateRegexp(value reflect.Value, check string) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		exp, err := regexp.Compile(check)
		if err != nil {
			return fmt.Errorf(msgIncorrectTag, errInternal, checkRegexp, check)
		}
		if !exp.MatchString(value.String()) {
			return fmt.Errorf(msgWrongRegexp, check)
		}
	default:
		return fmt.Errorf(msgUnexpectedType, errInternal, checkRegexp)
	}

	return nil
}

// validateIn проверяет на вхождение значения в список.
func validateIn(value reflect.Value, check string) error {
	values := strings.Split(check, validateListDelim)
	if len(values) == 0 {
		return fmt.Errorf(msgIncorrectTag, errInternal, checkIn, check)
	}
	equal := false
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		for idx := range values {
			if value.String() == values[idx] {
				equal = true
				break
			}
		}
		if !equal {
			return fmt.Errorf(msgWrongIn, value.String(), check)
		}
	case reflect.Int:
		for idx := range values {
			valueInt, err := strconv.Atoi(values[idx])
			if err != nil {
				return fmt.Errorf(msgIncorrectTag, errInternal, checkIn, check)
			}
			if int(value.Int()) == valueInt {
				equal = true
				break
			}
		}
		if !equal {
			return fmt.Errorf(msgWrongIn, value.Int(), check)
		}
	default:
		return fmt.Errorf(msgUnexpectedType, errInternal, checkIn)
	}

	return nil
}

// validateMin проверяет на соответствие минимальному значению границы.
func validateMin(value reflect.Value, check string) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.Int:
		min, err := strconv.Atoi(check)
		if err != nil {
			return fmt.Errorf(msgIncorrectTag, errInternal, checkMin, check)
		}
		if int(value.Int()) < min {
			return fmt.Errorf(msgLessMin, int(value.Int()), min)
		}
	default:
		return fmt.Errorf(msgUnexpectedType, errInternal, checkMin)
	}

	return nil
}

// validateMax проверяет на соответствие максимальному значению.
func validateMax(value reflect.Value, check string) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.Int:
		max, err := strconv.Atoi(check)
		if err != nil {
			return fmt.Errorf(msgIncorrectTag, errInternal, checkMax, check)
		}
		if int(value.Int()) > max {
			return fmt.Errorf(msgGreaterMax, int(value.Int()), max)
		}
	default:
		return fmt.Errorf(msgUnexpectedType, errInternal, checkMax)
	}

	return nil
}

// Validate валидирует публичные поля входной структуры.
func Validate(v interface{}) error {
	validator := Validator{}
	err := validator.Struct(v)
	if err != nil {
		return fmt.Errorf("%w:\n%v", errSystem, err.Error())
	}
	return validator.Errors()
}
