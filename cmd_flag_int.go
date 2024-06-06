package cling

type intCmdInput struct {
	name         string
	defaultValue *int
	required     bool
	description  string
	lDescription string
	envs         []string
	validator    validatorAny
}

func NewIntCmdInput(name string) CmdInputWithDefaultAndValidator[int] {
	return &intCmdInput{
		name: name,
	}
}

func (f *intCmdInput) FromEnv(sources []string) CmdFlag {
	f.envs = sources
	return f
}

func (f *intCmdInput) envSources() []string {
	return f.envs
}

func (f *intCmdInput) Description() string {
	return f.description
}

func (f *intCmdInput) WithDescription(value string) CmdInput {
	f.description = value
	return f
}

func (f *intCmdInput) longDescription() string {
	return f.lDescription
}

func (f *intCmdInput) WithLongDescription(value string) CmdArg {
	f.lDescription = value
	return f
}

func (f *intCmdInput) Required() CmdInput {
	f.required = true
	return f
}

func (f *intCmdInput) WithDefault(value int) CmdInputWithDefaultAndValidator[int] {
	if f.defaultValue == nil {
		f.defaultValue = new(int)
	}
	*f.defaultValue = value
	return f
}

func (f *intCmdInput) Name() string {
	return f.name
}

func (f *intCmdInput) WithValidator(validator Validator[int]) CmdInputWithDefaultAndValidator[int] {
	f.validator = &genericValidatorWrapper[int]{validator: validator}
	return f
}

func (f *intCmdInput) getValidator() validatorAny {
	return f.validator
}

func (f *intCmdInput) AsFlag() CmdFlag {
	return f
}

func (f *intCmdInput) AsArgument() CmdArg {
	return f
}

func (f *intCmdInput) isRequired() bool {
	return f.required
}

func (f *intCmdInput) hasDefault() bool {
	return f.defaultValue != nil
}

func (f *intCmdInput) getDefault() any {
	if f.defaultValue != nil {
		return *f.defaultValue
	}
	return 0
}
