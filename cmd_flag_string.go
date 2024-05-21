package cling

type StringCmdInput struct {
	name         string
	description  string
	lDescription string
	defaultValue *string
	required     bool
	envs         []string

	validators []Validator[string]
}

func NewStringCmdInput(name string) CmdInputWithDefaultAndValidator[string] {
	return &StringCmdInput{
		name: name,
	}
}

func (f *StringCmdInput) FromEnv(sources []string) CmdFlag {
	f.envs = sources
	return f
}

func (f *StringCmdInput) envSources() []string {
	return f.envs
}

func (f *StringCmdInput) Required() CmdInput {
	f.required = true
	return f
}

func (f *StringCmdInput) WithDefault(value string) CmdInputWithDefaultAndValidator[string] {
	if f.defaultValue == nil {
		f.defaultValue = new(string)
	}
	*f.defaultValue = value
	return f
}

func (f *StringCmdInput) Description() string {
	return f.description
}

func (f *StringCmdInput) WithDescription(value string) CmdInput {
	f.description = value
	return f
}

func (f *StringCmdInput) Name() string {
	return f.name
}

func (f *StringCmdInput) WithValidators(validators ...Validator[string]) CmdInputWithDefaultAndValidator[string] {
	f.validators = validators
	return f
}

func (f *StringCmdInput) getValidators() []Validator[string] {
	return f.validators
}

func (f *StringCmdInput) AsFlag() CmdFlag {
	return f
}

func (f *StringCmdInput) AsArgument() CmdArg {
	return f
}

func (f *StringCmdInput) WithLongDescription(longDescription string) CmdArg {
	f.lDescription = longDescription
	return f
}

func (f *StringCmdInput) longDescription() string {
	return f.lDescription
}

func (f *StringCmdInput) hasDefault() bool {
	return f.defaultValue != nil
}

func (f *StringCmdInput) getDefault() any {
	if f.defaultValue != nil {
		return *f.defaultValue
	}
	return ""
}

func (f *StringCmdInput) isRequired() bool {
	return f.required
}
