package cling

import (
	"errors"
	"fmt"
)

var ErrValidatorFailed = errors.New("validation failed")

// Validator interface with non-generic Validate method
type Validator[T any] interface {
	Validate(value T) error
}

type validatorAny interface {
	Validate(value any) error
}

type genericValidatorWrapper[S any] struct {
	validator Validator[S]
}

func (g *genericValidatorWrapper[S]) Validate(value any) error {
	val, ok := value.(S)
	if !ok {
		return fmt.Errorf("invalid type: expected %T, got %T", new(S), value)
	}
	return g.validator.Validate(val)
}

// NoOpValidator returns a no-op validator for any type
func NoOpValidator() validatorAny {
	return &noOpValidator{}
}

type noOpValidator struct{}

func (v *noOpValidator) Validate(value any) error {
	return nil
}

// ComposeValidator composes multiple validators for a specific type
func ComposeValidator[T any](validators ...Validator[T]) Validator[T] {
	return &compositeValidator[T]{validators: validators}
}

type compositeValidator[T any] struct {
	validators []Validator[T]
}

func (v *compositeValidator[T]) Validate(value T) error {
	for _, validator := range v.validators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	return nil
}
