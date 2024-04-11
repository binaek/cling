package cling

import "context"

type ClingContextKey string

const (
	ContextKeyCommand ClingContextKey = "command"
)

func WithCommand(ctx context.Context, command *Command) context.Context {
	return context.WithValue(ctx, ContextKeyCommand, command)
}

func CommandFromContext(ctx context.Context) (*Command, bool) {
	command, ok := ctx.Value(ContextKeyCommand).(*Command)
	return command, ok
}
