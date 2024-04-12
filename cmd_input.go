package cling

type CmdInput interface {
	Name() string
	Description() string
	Required() CmdInput
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
}

type CmdFlag interface {
	CmdInput
}

type CmdArg interface {
	CmdInput
	WithLongDescription(string) CmdArg
	LongDescription() string
}
