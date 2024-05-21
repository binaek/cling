package cling

import (
	"fmt"
	"reflect"
	"strconv"
)

func setFieldFromString(field reflect.Value, valueStr string, validators ...Validator[any]) error {
	var value any
	var err error

	switch field.Kind() {
	case reflect.String:
		value, err = parseString(valueStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err = parseInt(valueStr, field.Type().Bits())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err = parseUint(valueStr, field.Type().Bits())
	case reflect.Bool:
		value, err = parseBool(valueStr)
	case reflect.Float32, reflect.Float64:
		value, err = parseFloat(valueStr, field.Type().Bits())
	case reflect.Slice:
		return setSlice(field, valueStr, validators...)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type().Kind())
	}

	if err != nil {
		return err
	}

	if err := runValidators(value, validators); err != nil {
		return err
	}

	field.Set(reflect.ValueOf(value).Convert(field.Type()))
	return nil
}

func setSlice(field reflect.Value, valueStr string, validators ...Validator[any]) error {
	elemType := field.Type().Elem()
	var value any
	var err error

	switch elemType.Kind() {
	case reflect.String:
		value, err = parseString(valueStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err = parseInt(valueStr, elemType.Bits())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err = parseUint(valueStr, elemType.Bits())
	case reflect.Bool:
		value, err = parseBool(valueStr)
	case reflect.Float32, reflect.Float64:
		value, err = parseFloat(valueStr, elemType.Bits())
	default:
		return fmt.Errorf("unsupported slice element type: %s", elemType.Kind())
	}

	if err != nil {
		return err
	}

	if err := runValidators(value, validators); err != nil {
		return err
	}

	field.Set(reflect.Append(field, reflect.ValueOf(value).Convert(elemType)))
	return nil
}

func parseString(valueStr string) (string, error) {
	return valueStr, nil
}

func parseInt(valueStr string, bits int) (int64, error) {
	return strconv.ParseInt(valueStr, 10, bits)
}

func parseUint(valueStr string, bits int) (uint64, error) {
	return strconv.ParseUint(valueStr, 10, bits)
}

func parseBool(valueStr string) (bool, error) {
	return strconv.ParseBool(valueStr)
}

func parseFloat(valueStr string, bits int) (float64, error) {
	return strconv.ParseFloat(valueStr, bits)
}

func runValidators(value any, validators []Validator[any]) error {
	for _, validator := range validators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	return nil
}
