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
	taskLen    = "len"
	taskRegexp = "regexp"
	taskIn     = "in"
	taskMin    = "min"
	taskMax    = "max"
)

var (
	msgExpectedStruct = "expected struct but returned %s value"
	msgExpectedValue  = "%w: tag '%v' have not value"
	msgUndefinedTag   = "%w: tag '%v' is undefined"
	msgIncorrectTag   = "%w: incorrect value type for tag '%v' value '%v'"
	msgUnexpectedType = "%w: unexpected type for validation '%v' function"

	msgWrongLen    = "wrong length: must be '%d' not a '%d'"
	msgWrongRegexp = "wrong regexp: value does not math pattern '%v'"
	msgErrMin      = "value is less than the limit, must be greater than '%d'"
	msgErrMax      = "value is greater than the limit, must be less than '%d'"
	msgWrongIn     = "wrong in: value '%v' does not math any values from list '%v'"

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

// Validator объект, валидирующий структуру.
type Validator struct {
	validationErrs ValidationErrors // Ошибки валидации
}

// Field валидирует атомарное поле.
func (v *Validator) Field(name string, value reflect.Value, tags string) error {
	checks := strings.Split(tags, validateTagDelim)
	for _, check := range checks {
		var err error
		metadata := strings.Split(check, validateCheckDelim)
		if len(metadata) != 2 {
			return fmt.Errorf(msgExpectedValue, errInternal, metadata[0])
		}

		nameCheck, valueCheck := metadata[0], metadata[1]
		switch nameCheck {
		case taskLen:
			err = validateEqualLen(value, valueCheck)
		case taskRegexp:
			err = validateRegexp(value, valueCheck)
		case taskIn:
			err = validateIn(value, valueCheck)
		case taskMin:
			err = validateMin(value, valueCheck)
		case taskMax:
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

// validateStruct валидирует структуру.
func (v *Validator) validateStruct(value reflect.Value, valueType reflect.Type) error {
	for i := 0; i < valueType.NumField(); i++ {
		field := value.Field(i)
		fieldType := valueType.Field(i)

		tags, ok := fieldType.Tag.Lookup(validateTag)
		if !ok {
			continue
		}

		fieldName := valueType.Field(i).Name
		var err error

		switch {
		case field.Kind() == reflect.Struct && isNested(tags):
			err = v.validateStruct(field, field.Type())
			if err != nil {
				return err
			}
		case field.Kind() == reflect.Slice && isNested(tags):
			for j := 0; j < field.Len(); j++ {
				err = v.validateStruct(field.Index(j), field.Index(j).Type())
				if err != nil {
					return err
				}
			}
		case field.Kind() != reflect.Struct && field.Kind() != reflect.Slice:
			err = v.Field(fieldName, field, tags)
			if err != nil {
				return err
			}
		case field.Kind() != reflect.Struct && field.Kind() == reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				err = v.Field(fieldName, field.Index(j), tags)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Struct валидирует переданную структуру.
func (v *Validator) Struct(i interface{}) error {
	if v == nil {
		return fmt.Errorf(msgExpectedStruct, "nil")
	}

	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Struct {
		return fmt.Errorf(msgExpectedStruct, value.Kind())
	}
	return v.validateStruct(value, value.Type())
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
			return fmt.Errorf(msgIncorrectTag, errInternal, taskLen, value)
		}
		if value.Len() != long {
			return fmt.Errorf(msgWrongLen, long, value.Len())
		}
	default:
		return fmt.Errorf(msgUnexpectedType, errInternal, taskLen)
	}

	return nil
}

// validateRegexp проверяет на соответствие регулярному выражению.
func validateRegexp(value reflect.Value, check string) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		exp, err := regexp.Compile(check)
		if err != nil {
			return fmt.Errorf(msgIncorrectTag, errInternal, taskRegexp, check)
		}
		if !exp.MatchString(value.String()) {
			return fmt.Errorf(msgWrongRegexp, check)
		}
	default:
		return fmt.Errorf(msgUnexpectedType, errInternal, taskRegexp)
	}

	return nil
}

// validateIn проверяет на вхождение значения в список.
func validateIn(value reflect.Value, check string) error {
	values := strings.Split(check, validateListDelim)
	if len(values) == 0 {
		return fmt.Errorf(msgIncorrectTag, errInternal, taskIn, check)
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
				return fmt.Errorf(msgIncorrectTag, errInternal, taskIn, check)
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
		return fmt.Errorf(msgUnexpectedType, errInternal, taskIn)
	}

	return nil
}

// validateMin проверяет на соответствие минимальному значению границы.
func validateMin(value reflect.Value, check string) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.Int:
		min, err := strconv.Atoi(check)
		if err != nil {
			return fmt.Errorf(msgIncorrectTag, errInternal, taskMin, check)
		}
		if int(value.Int()) < min {
			return fmt.Errorf(msgErrMin, min)
		}
	default:
		return fmt.Errorf(msgUnexpectedType, errInternal, taskMin)
	}

	return nil
}

// validateMax проверяет на соответствие максимальному значению.
func validateMax(value reflect.Value, check string) error {
	switch value.Kind() { //nolint:exhaustive
	case reflect.Int:
		max, err := strconv.Atoi(check)
		if err != nil {
			return fmt.Errorf(msgIncorrectTag, errInternal, taskMax, check)
		}
		if int(value.Int()) > max {
			return fmt.Errorf(msgErrMax, max)
		}
	default:
		return fmt.Errorf(msgUnexpectedType, errInternal, taskMax)
	}

	return nil
}

// isNested возвращает признак, что вложенная структура требует валидации.
func isNested(tags string) bool {
	return tags == nestedTag
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
