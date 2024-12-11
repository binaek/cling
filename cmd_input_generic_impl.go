package cling

// NewIntCmdInput creates a new integer command input with the given name.
func NewIntCmdInput(name string) CmdInputWithDefaultAndValidator[int] {
	return newGenericCmdInput[int](name)
}

// NewStringCmdInput creates a new string command input with the given name.
func NewStringCmdInput(name string) CmdInputWithDefaultAndValidator[string] {
	return newGenericCmdInput[string](name)
}

// NewBoolCmdInput creates a new boolean command input with the given name.
func NewBoolCmdInput(name string) CmdInputWithDefaultAndValidator[bool] {
	return newGenericCmdInput[bool](name)
}

type genericCmdInput[T int | string | bool] struct {
	name         string
	defaultValue *T
	required     bool
	description  string
	lDescription string
	envs         []string
	validator    validatorAny
}

func newGenericCmdInput[T int | string | bool](name string) CmdInputWithDefaultAndValidator[T] {
	return &genericCmdInput[T]{
		name: name,
	}
}

func (f *genericCmdInput[T]) FromEnv(sources []string) CmdFlag {
	f.envs = sources
	return f
}

func (f *genericCmdInput[T]) envSources() []string {
	return f.envs
}

func (f *genericCmdInput[T]) Description() string {
	return f.description
}

func (f *genericCmdInput[T]) WithDescription(value string) CmdInput {
	f.description = value
	return f
}

func (f *genericCmdInput[T]) longDescription() string {
	return f.lDescription
}

func (f *genericCmdInput[T]) WithLongDescription(value string) CmdArg {
	f.lDescription = value
	return f
}

func (f *genericCmdInput[T]) Required() CmdInput {
	f.required = true
	return f
}

func (f *genericCmdInput[T]) WithDefault(value T) CmdInputWithDefaultAndValidator[T] {
	if f.defaultValue == nil {
		f.defaultValue = new(T)
	}
	*f.defaultValue = value
	return f
}

func (f *genericCmdInput[T]) Name() string {
	return f.name
}

func (f *genericCmdInput[T]) WithValidator(validator Validator[T]) CmdInputWithDefaultAndValidator[T] {
	f.validator = &genericValidatorWrapper[T]{validator: validator}
	return f
}

func (f *genericCmdInput[T]) getValidator() validatorAny {
	return f.validator
}

func (f *genericCmdInput[T]) AsFlag() CmdFlag {
	return f
}

func (f *genericCmdInput[T]) AsArgument() CmdArg {
	return f
}

func (f *genericCmdInput[T]) isRequired() bool {
	return f.required
}

func (f *genericCmdInput[T]) hasDefault() bool {
	return f.defaultValue != nil
}

func (f *genericCmdInput[T]) getDefault() any {
	if f.defaultValue != nil {
		return *f.defaultValue
	}
	return 0
}
