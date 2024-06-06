package cling

type cmdInputGenericSlice[T comparable] struct {
	name         string
	description  string
	lDescription string
	defaultValue []T
	required     bool
	envs         []string
	validator    validatorAny
}

func NewCmdSliceInput[T comparable](name string) CmdInputWithDefaultAndValidator[[]T] {
	return &cmdInputGenericSlice[T]{
		name:         name,
		defaultValue: nil,
	}
}

func (f *cmdInputGenericSlice[T]) AsArgument() CmdArg {
	return f
}

func (f *cmdInputGenericSlice[T]) AsFlag() CmdFlag {
	return f
}

func (f *cmdInputGenericSlice[T]) WithDescription(value string) CmdInput {
	f.description = value
	return f
}

func (f *cmdInputGenericSlice[T]) Description() string {
	return f.description
}

func (f *cmdInputGenericSlice[T]) Name() string {
	return f.name
}

func (f *cmdInputGenericSlice[T]) Required() CmdInput {
	f.required = true
	return f
}

func (f *cmdInputGenericSlice[T]) WithDefault(value []T) CmdInputWithDefaultAndValidator[[]T] {
	f.defaultValue = value
	return f
}

func (f *cmdInputGenericSlice[T]) WithValidator(validator Validator[[]T]) CmdInputWithDefaultAndValidator[[]T] {
	f.validator = &genericValidatorWrapper[[]T]{validator: validator}
	return f
}

func (f *cmdInputGenericSlice[T]) FromEnv(sources []string) CmdFlag {
	f.envs = sources
	return f
}

func (f *cmdInputGenericSlice[T]) WithLongDescription(value string) CmdArg {
	f.lDescription = value
	return f
}

func (f *cmdInputGenericSlice[T]) getValidator() validatorAny {
	return f.validator
}

func (f *cmdInputGenericSlice[T]) envSources() []string {
	return f.envs
}

func (f *cmdInputGenericSlice[T]) longDescription() string {
	return f.lDescription
}

func (f *cmdInputGenericSlice[T]) hasDefault() bool {
	return f.defaultValue != nil
}

func (f *cmdInputGenericSlice[T]) getDefault() any {
	return f.defaultValue
}

func (f *cmdInputGenericSlice[T]) isRequired() bool {
	return f.required
}
