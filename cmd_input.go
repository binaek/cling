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

type ValidatorProvider interface {
	getValidator() validatorAny
}

type CmdInputWithDefaultAndValidator[S any] interface {
	CmdInput
	ValidatorProvider
	WithDefault(value S) CmdInputWithDefaultAndValidator[S]
	WithValidator(validator Validator[S]) CmdInputWithDefaultAndValidator[S]
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
