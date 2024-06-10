package cling

import (
	"github.com/pkg/errors"
)

type stringValidator struct {
	fn func(value string) error
}

func (v *stringValidator) Validate(value string) error {
	return v.fn(value)
}

var ErrStringLen = errors.New("string length is not in range")

func NewStringLengthValidator(min, max int) Validator[string] {
	return &stringValidator{
		fn: func(value string) error {
			if len(value) < min || len(value) > max {
				return errors.Wrapf(ErrStringLen, "value %s is not in range [%d, %d]", value, min, max)
			}
			return nil
		},
	}
}
