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

const (
	msgExpectedStruct = "expected struct but returned %s value"
	msgIncorrectTag   = "%w: incorrect value type for tag '%v' value '%v'"
	msgUnexpectedType = "%w: unexpected type for validation '%v' function"
	msgWrongIn        = "value '%v' does not math any values from list '%v'"
)

var (
	errInternal   = errors.New("internal error")
	errSystem     = errors.New("system error")
	errValidation = errors.New("validation error")
)

// ValidationError ошибки валидации.
type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result strings.Builder
	for i := range v {
		result.WriteString(fmt.Sprintf("%v: %v\n", v[i].Field, v[i].Err))
	}
	return result.String()
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
		if fieldStruct.Anonymous {
			continue
		}

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
		err = v.validateStruct(field, field.Type())
	case field.Kind() == reflect.Slice && isNested(tags):
		for j := 0; j < field.Len(); j++ {
			if err = v.validateStruct(field.Index(j), field.Index(j).Type()); err != nil {
				break
			}
		}
	case field.Kind() != reflect.Struct && field.Kind() != reflect.Slice:
		err = v.validateField(name, field, tags)
	case field.Kind() != reflect.Struct && field.Kind() == reflect.Slice:
		for j := 0; j < field.Len(); j++ {
			if err = v.validateField(name, field.Index(j), tags); err != nil {
				break
			}
		}
	}
	return err
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
			return fmt.Errorf("%w: tag '%v' have not value", errInternal, metadata[0])
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
			err = fmt.Errorf("%w: tag '%v' is undefined", errInternal, metadata[0])
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
func (v *Validator) Errors() ValidationErrors {
	return v.validationErrs
}

// validateEqualLen проверяет равенство длины.
func validateEqualLen(value reflect.Value, check string) error {
	if value.Kind() != reflect.String {
		return fmt.Errorf(msgUnexpectedType, errInternal, checkLen)
	}
	long, err := strconv.Atoi(check)
	if err != nil {
		return fmt.Errorf(msgIncorrectTag, errInternal, checkLen, value)
	}
	if value.Len() != long {
		return fmt.Errorf("the length of the value must be '%d' not a '%d'", long, value.Len())
	}

	return nil
}

// validateRegexp проверяет на соответствие регулярному выражению.
func validateRegexp(value reflect.Value, check string) error {
	if value.Kind() != reflect.String {
		return fmt.Errorf(msgUnexpectedType, errInternal, checkRegexp)
	}

	exp, err := regexp.Compile(check)
	if err != nil {
		return fmt.Errorf(msgIncorrectTag, errInternal, checkRegexp, check)
	}
	if !exp.MatchString(value.String()) {
		return fmt.Errorf("value does not math pattern '%v'", check)
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
	if value.Kind() != reflect.Int {
		return fmt.Errorf(msgUnexpectedType, errInternal, checkMin)
	}
	min, err := strconv.Atoi(check)
	if err != nil {
		return fmt.Errorf(msgIncorrectTag, errInternal, checkMin, check)
	}
	if int(value.Int()) < min {
		return fmt.Errorf("value '%d' is less than the limit, must be greater than '%d'", int(value.Int()), min)
	}
	return nil
}

// validateMax проверяет на соответствие максимальному значению.
func validateMax(value reflect.Value, check string) error {
	if value.Kind() != reflect.Int {
		return fmt.Errorf(msgUnexpectedType, errInternal, checkMax)
	}
	max, err := strconv.Atoi(check)
	if err != nil {
		return fmt.Errorf(msgIncorrectTag, errInternal, checkMax, check)
	}
	if int(value.Int()) > max {
		return fmt.Errorf("value '%d' is greater than the limit, must be less than '%d'", int(value.Int()), max)
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
	validationErrs := validator.Errors()

	if len(validationErrs) == 0 {
		return nil
	}

	var result strings.Builder
	for i := range validationErrs {
		result.WriteString(fmt.Sprintf("%v: %v\n", validationErrs[i].Field, validationErrs[i].Err))
	}
	return fmt.Errorf("%w:\n%v", errValidation, result.String())
}
