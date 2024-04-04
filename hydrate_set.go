package cling

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func setFieldFromString(field reflect.Value, valueStr string) error {
	switch field.Kind() {
	case reflect.String:
		return setString(field, valueStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setInt(field, valueStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUint(field, valueStr)
	case reflect.Bool:
		return setBool(field, valueStr)
	case reflect.Float32, reflect.Float64:
		return setFloat(field, valueStr)
	case reflect.Slice:
		return setSlice(field, valueStr)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type().Kind())
	}
}
func setSlice(field reflect.Value, valueStr string) error {
	elemType := field.Type().Elem()
	switch elemType.Kind() {
	case reflect.String:
		return appendString(field, valueStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return appendInt(field, valueStr, elemType)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return appendUint(field, valueStr, elemType)
	case reflect.Bool:
		return appendBool(field, valueStr)
	case reflect.Float32, reflect.Float64:
		return appendFloat(field, valueStr, elemType)
	default:
		return fmt.Errorf("unsupported slice element type: %s", elemType.Kind())
	}
}

func setString(field reflect.Value, valueStr string) error {
	field.SetString(valueStr)
	return nil
}

func setInt(field reflect.Value, valueStr string) error {
	valueInt, err := strconv.ParseInt(valueStr, 10, field.Type().Bits())
	if err != nil {
		return err
	}
	field.SetInt(valueInt)
	return nil
}

func setUint(field reflect.Value, valueStr string) error {
	valueUint, err := strconv.ParseUint(valueStr, 10, field.Type().Bits())
	if err != nil {
		return err
	}
	field.SetUint(valueUint)
	return nil
}

func setBool(field reflect.Value, valueStr string) error {
	valueBool, err := strconv.ParseBool(valueStr)
	if err != nil {
		if !errors.Is(err, strconv.ErrSyntax) {
			return err
		}
		valueBool = true
	}
	field.SetBool(valueBool)
	return nil
}

func setFloat(field reflect.Value, valueStr string) error {
	valueFloat, err := strconv.ParseFloat(valueStr, field.Type().Bits())
	if err != nil {
		return err
	}
	field.SetFloat(valueFloat)
	return nil
}

func appendString(field reflect.Value, valueStr string) error {
	field.Set(reflect.Append(field, reflect.ValueOf(valueStr)))
	return nil
}

func appendInt(field reflect.Value, valueStr string, elemType reflect.Type) error {
	valueInt, err := strconv.ParseInt(valueStr, 10, elemType.Bits())
	if err != nil {
		return err
	}
	field.Set(reflect.Append(field, reflect.ValueOf(valueInt).Convert(elemType)))
	return nil
}

func appendUint(field reflect.Value, valueStr string, elemType reflect.Type) error {
	valueUint, err := strconv.ParseUint(valueStr, 10, elemType.Bits())
	if err != nil {
		return err
	}
	field.Set(reflect.Append(field, reflect.ValueOf(valueUint).Convert(elemType)))
	return nil
}

func appendBool(field reflect.Value, valueStr string) error {
	valueBool, err := strconv.ParseBool(valueStr)
	if err != nil {
		if !errors.Is(err, strconv.ErrSyntax) {
			return err
		}
		valueBool = true
	}
	field.Set(reflect.Append(field, reflect.ValueOf(valueBool)))
	return nil
}

func appendFloat(field reflect.Value, valueStr string, elemType reflect.Type) error {
	valueFloat, err := strconv.ParseFloat(valueStr, elemType.Bits())
	if err != nil {
		return err
	}
	field.Set(reflect.Append(field, reflect.ValueOf(valueFloat).Convert(elemType)))
	return nil
}
