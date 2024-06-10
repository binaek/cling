package cling

import (
	"github.com/pkg/errors"
)

type intValidator struct {
	fn func(value int) error
}

func (v *intValidator) Validate(value int) error {
	return v.fn(value)
}

func NewIntRangeValidator(min, max int) Validator[int] {
	return &intValidator{
		fn: func(value int) error {
			if value < min || value > max {
				return errors.Wrapf(ErrValidatorFailed, "value %d is not in range [%d, %d]", value, min, max)
			}
			return nil
		},
	}
}
