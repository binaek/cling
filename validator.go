package cling

type Validator[T any] interface {
	Validate(value T) error
}
