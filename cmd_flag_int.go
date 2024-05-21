package cling

type IntCmdInput struct {
	name         string
	defaultValue *int
	required     bool
	description  string
	lDescription string
	envs         []string

	validators []Validator[int]
}

func NewIntCmdInput(name string) CmdInputWithDefaultAndValidator[int] {
	return &IntCmdInput{
		name: name,
	}
}

func (f *IntCmdInput) FromEnv(sources []string) CmdFlag {
	f.envs = sources
	return f
}

func (f *IntCmdInput) envSources() []string {
	return f.envs
}

func (f *IntCmdInput) Description() string {
	return f.description
}

func (f *IntCmdInput) WithDescription(value string) CmdInput {
	f.description = value
	return f
}

func (f *IntCmdInput) longDescription() string {
	return f.lDescription
}

func (f *IntCmdInput) WithLongDescription(value string) CmdArg {
	f.lDescription = value
	return f
}

func (f *IntCmdInput) Required() CmdInput {
	f.required = true
	return f
}

func (f *IntCmdInput) WithDefault(value int) CmdInputWithDefaultAndValidator[int] {
	if f.defaultValue == nil {
		f.defaultValue = new(int)
	}
	*f.defaultValue = value
	return f
}

func (f *IntCmdInput) Name() string {
	return f.name
}

func (f *IntCmdInput) WithValidators(validators ...Validator[int]) CmdInputWithDefaultAndValidator[int] {
	f.validators = validators
	return f
}

func (f *IntCmdInput) getValidators() []Validator[int] {
	return f.validators
}

func (f *IntCmdInput) AsFlag() CmdFlag {
	return f
}

func (f *IntCmdInput) AsArgument() CmdArg {
	return f
}

func (f *IntCmdInput) isRequired() bool {
	return f.required
}

func (f *IntCmdInput) hasDefault() bool {
	return f.defaultValue != nil
}

func (f *IntCmdInput) getDefault() any {
	if f.defaultValue != nil {
		return *f.defaultValue
	}
	return 0
}
