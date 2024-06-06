package cling

type CmdInput interface {
	// Name returns the name of the command input.
	Name() string
	// Required marks the command input as required.
	Required() CmdInput
	// WithDescription sets the description of the command input.
	WithDescription(description string) CmdInput
	// Description returns the description of the command input.
	Description() string
	// AsFlag returns the command input as a flag.
	AsFlag() CmdFlag
	// AsArgument returns the command input as an argument.
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
	// WithDefault sets the default value of the command input.
	WithDefault(value S) CmdInputWithDefaultAndValidator[S]
	// WithValidator sets the validator of the command input.
	WithValidator(validator Validator[S]) CmdInputWithDefaultAndValidator[S]
}

type CmdFlag interface {
	CmdInput
	// FromEnv sets the environment sources of the command flag.
	// The environment sources are used to get the default value of the flag.
	// The default value is the first non-empty value found in the environment.
	// This will override the default value set by WithDefault.
	FromEnv([]string) CmdFlag
	envSources() []string
}

type CmdArg interface {
	CmdInput
	// WithLongDescription sets the long description of the command argument.
	WithLongDescription(string) CmdArg
	longDescription() string
}
