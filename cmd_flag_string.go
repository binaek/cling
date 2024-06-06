package cling

type stringCmdInput struct {
	name         string
	description  string
	lDescription string
	defaultValue *string
	required     bool
	envs         []string
	validator    validatorAny
}

func NewStringCmdInput(name string) CmdInputWithDefaultAndValidator[string] {
	return &stringCmdInput{
		name: name,
	}
}

func (f *stringCmdInput) FromEnv(sources []string) CmdFlag {
	f.envs = sources
	return f
}

func (f *stringCmdInput) envSources() []string {
	return f.envs
}

func (f *stringCmdInput) Required() CmdInput {
	f.required = true
	return f
}

func (f *stringCmdInput) WithDefault(value string) CmdInputWithDefaultAndValidator[string] {
	if f.defaultValue == nil {
		f.defaultValue = new(string)
	}
	*f.defaultValue = value
	return f
}

func (f *stringCmdInput) Description() string {
	return f.description
}

func (f *stringCmdInput) WithDescription(value string) CmdInput {
	f.description = value
	return f
}

func (f *stringCmdInput) Name() string {
	return f.name
}

func (f *stringCmdInput) WithValidator(validator Validator[string]) CmdInputWithDefaultAndValidator[string] {
	f.validator = &genericValidatorWrapper[string]{validator: validator}
	return f
}

func (f *stringCmdInput) getValidator() validatorAny {
	return f.validator
}
func (f *stringCmdInput) AsFlag() CmdFlag {
	return f
}

func (f *stringCmdInput) AsArgument() CmdArg {
	return f
}

func (f *stringCmdInput) WithLongDescription(longDescription string) CmdArg {
	f.lDescription = longDescription
	return f
}

func (f *stringCmdInput) longDescription() string {
	return f.lDescription
}

func (f *stringCmdInput) hasDefault() bool {
	return f.defaultValue != nil
}

func (f *stringCmdInput) getDefault() any {
	if f.defaultValue != nil {
		return *f.defaultValue
	}
	return ""
}

func (f *stringCmdInput) isRequired() bool {
	return f.required
}
