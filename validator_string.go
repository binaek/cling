package cling

import "fmt"

type StringLengthValidator struct {
	min int
	max int
}

func NewStringLengthValidator(min, max int) *StringLengthValidator {
	return &StringLengthValidator{
		min: min,
		max: max,
	}
}

func (v *StringLengthValidator) Validate(value string) error {
	if len(value) < v.min || len(value) > v.max {
		return fmt.Errorf("value %s length is out of range [%d, %d]", value, v.min, v.max)
	}
	return nil
}
