package cling

import "context"

type ClingContextKey string

const (
	ContextKeyCommand ClingContextKey = "command"
)

func contextWithCommand(ctx context.Context, command *Command) context.Context {
	return context.WithValue(ctx, ContextKeyCommand, command)
}

func commandFromContext(ctx context.Context) (*Command, bool) {
	command, ok := ctx.Value(ContextKeyCommand).(*Command)
	return command, ok
}
