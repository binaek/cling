package cling

import "fmt"

type EnumValidator[T comparable] struct {
	allowedValues []T
}

func NewEnumValidator[T comparable](allowedValues ...T) *EnumValidator[T] {
	return &EnumValidator[T]{
		allowedValues: allowedValues,
	}
}

func (v *EnumValidator[T]) Validate(value T) error {
	for _, allowedValue := range v.allowedValues {
		if value == allowedValue {
			return nil
		}
	}
	return fmt.Errorf("value %v is not one of %v", value, v.allowedValues)
}
