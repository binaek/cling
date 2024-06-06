package cling

type cmdInputSlice[T comparable] struct {
	name         string
	description  string
	lDescription string
	defaultValue []T
	required     bool
	envs         []string
	validator    validatorAny
}

func NewCmdSliceInput[T comparable](name string) CmdInputWithDefaultAndValidator[[]T] {
	return &cmdInputSlice[T]{
		name: name,
	}
}

func (f *cmdInputSlice[T]) AsArgument() CmdArg {
	return f
}

func (f *cmdInputSlice[T]) AsFlag() CmdFlag {
	return f
}

func (f *cmdInputSlice[T]) WithDescription(value string) CmdInput {
	f.description = value
	return f
}

func (f *cmdInputSlice[T]) Description() string {
	return f.description
}

func (f *cmdInputSlice[T]) Name() string {
	return f.name
}

func (f *cmdInputSlice[T]) Required() CmdInput {
	f.required = true
	return f
}

func (f *cmdInputSlice[T]) WithDefault(value []T) CmdInputWithDefaultAndValidator[[]T] {
	f.defaultValue = value
	return f
}

func (f *cmdInputSlice[T]) WithValidator(validator Validator[[]T]) CmdInputWithDefaultAndValidator[[]T] {
	f.validator = &genericValidatorWrapper[[]T]{validator: validator}
	return f
}

func (f *cmdInputSlice[T]) FromEnv(sources []string) CmdFlag {
	f.envs = sources
	return f
}

func (f *cmdInputSlice[T]) WithLongDescription(value string) CmdArg {
	f.lDescription = value
	return f
}

func (f *cmdInputSlice[T]) getValidator() validatorAny {
	return f.validator
}

func (f *cmdInputSlice[T]) envSources() []string {
	return f.envs
}

func (f *cmdInputSlice[T]) longDescription() string {
	return f.lDescription
}

func (f *cmdInputSlice[T]) hasDefault() bool {
	return f.defaultValue != nil
}

func (f *cmdInputSlice[T]) getDefault() any {
	if f.defaultValue != nil {
		return f.defaultValue
	}
	return ""
}

func (f *cmdInputSlice[T]) isRequired() bool {
	return f.required
}
