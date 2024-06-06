package cling

import "github.com/pkg/errors"

type Comparator[T any] func(a T) error

type comparatorValidator[T any] struct {
	comparator Comparator[T]
}

// NewComparatorValidator creates a new validator that checks if the value satisfies the comparator.
func NewComparatorValidator[T any](comparator Comparator[T]) Validator[T] {
	return &comparatorValidator[T]{
		comparator: comparator,
	}
}

func (v *comparatorValidator[T]) Validate(value T) error {
	return v.comparator(value)
}

// NewEnumValidator creates a new validator that checks if the value is one of the allowed values.
func NewEnumValidator[T comparable](allowedValues ...T) Validator[T] {
	return NewComparatorValidator[T](func(value T) error {
		for _, allowed := range allowedValues {
			if value == allowed {
				return nil
			}
		}
		return errors.Wrapf(ErrValidatorFailed, "value '%v' is not in the allowed enum values", value)
	})
}
