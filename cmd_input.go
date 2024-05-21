package cling

type CmdInput interface {
	Name() string
	Description() string
	Required() CmdInput
	WithDescription(description string) CmdInput

	AsFlag() CmdFlag
	AsArgument() CmdArg

	isRequired() bool
	hasDefault() bool
	getDefault() any
}

type CmdInputWithDefaultAndValidator[S any] interface {
	CmdInput
	WithDefault(value S) CmdInputWithDefaultAndValidator[S]
	WithValidators(validators ...Validator[S]) CmdInputWithDefaultAndValidator[S]

	getValidators() []Validator[S]
}

type CmdFlag interface {
	CmdInput
	FromEnv([]string) CmdFlag
	envSources() []string
}

type CmdArg interface {
	CmdInput
	WithLongDescription(string) CmdArg
	longDescription() string
}
