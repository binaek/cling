package cling

import "fmt"

type IntRangeValidator struct {
	min int
	max int
}

func NewIntRangeValidator(min, max int) *IntRangeValidator {
	return &IntRangeValidator{
		min: min,
		max: max,
	}
}

func (v *IntRangeValidator) Validate(value int) error {
	if value < v.min || value > v.max {
		return fmt.Errorf("value %d is out of range [%d, %d]", value, v.min, v.max)
	}
	return nil
}
