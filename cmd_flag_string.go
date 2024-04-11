package cling

type StringCmdInput struct {
	name            string
	description     string
	longDescription string
	defaultValue    *string
	required        bool
}

func NewStringCmdInput(name string) CmdInputWithDefaultAndValidator[string] {
	return &StringCmdInput{
		name: name,
	}
}

func (f *StringCmdInput) Required() CmdInput {
	f.required = true
	return f
}

func (f *StringCmdInput) isRequired() bool {
	return f.required
}

func (f *StringCmdInput) WithDefault(value string) CmdInputWithDefaultAndValidator[string] {
	if f.defaultValue == nil {
		f.defaultValue = new(string)
	}
	*f.defaultValue = value
	return f
}

func (f *StringCmdInput) getDefault() any {
	if f.defaultValue != nil {
		return *f.defaultValue
	}
	return ""
}

func (f *StringCmdInput) hasDefault() bool {
	return f.defaultValue != nil
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
	return f
}

func (f *StringCmdInput) AsFlag() CmdFlag {
	return f
}

func (f *StringCmdInput) AsArgument() CmdArg {
	return f
}

func (f *StringCmdInput) WithLongDescription(longDescription string) CmdArg {
	f.longDescription = longDescription
	return f
}

func (f *StringCmdInput) LongDescription() string {
	return f.longDescription
}
