package cling

import (
	"context"
	"slices"
)

type CommandHook func(ctx context.Context, args []string) error
type CommandHandler func(ctx context.Context, args []string) error

func NoOpHook(ctx context.Context, args []string) error {
	return nil
}

type Command struct {
	name            string
	description     string
	longDescription string
	action          CommandHandler
	flags           map[string]CmdFlag
	arguments       []CmdArg
	children        []*Command
	parent          *Command

	// hooks
	preRun            CommandHook
	postRun           CommandHook
	persistentPreRun  CommandHook
	persistentPostRun CommandHook
}

func NewCommand(name string, action CommandHandler) *Command {
	command := &Command{
		name:              name,
		action:            action,
		flags:             make(map[string]CmdFlag),
		arguments:         []CmdArg{},
		preRun:            NoOpHook,
		postRun:           NoOpHook,
		persistentPreRun:  NoOpHook,
		persistentPostRun: NoOpHook,
	}
	return command
}

func (c *Command) WithPreRun(hook CommandHook) *Command {
	c.preRun = hook
	return c
}

func (c *Command) WithPostRun(hook CommandHook) *Command {
	c.postRun = hook
	return c
}

func (c *Command) WithPersistentPreRun(hook CommandHook) *Command {
	c.persistentPreRun = hook
	return c
}

func (c *Command) WithPersistentPostRun(hook CommandHook) *Command {
	c.persistentPostRun = hook
	return c
}

func (c *Command) WithDescription(description string) *Command {
	c.description = description
	return c
}

func (c *Command) WithLongDescription(longDescription string) *Command {
	c.longDescription = longDescription
	return c
}

func (command *Command) WithChildCommand(cmd *Command) *Command {
	command.children = append(command.children, cmd)
	cmd.parent = command
	return command
}

func (command *Command) WithFlag(flag CmdFlag) *Command {
	command.flags[flag.Name()] = flag
	return command
}

func (command *Command) WithArgument(arg CmdArg) *Command {
	command.arguments = append(command.arguments, arg)
	return command
}

func (command *Command) execute(ctx context.Context, args []string) error {
	if err := command.executePrerun(ctx, args); err != nil {
		return err
	}
	if err := command.action(ctx, args); err != nil {
		return err
	}
	if err := command.executePostRun(ctx, args); err != nil {
		return err
	}
	return nil
}

func (command *Command) executePrerun(ctx context.Context, args []string) error {
	path := command.pathToRoot()
	slices.Reverse(path)
	for _, cmd := range path {
		if cmd.persistentPreRun != nil {
			if err := cmd.persistentPreRun(ctx, args); err != nil {
				return err
			}
		}
	}
	if command.preRun != nil {
		if err := command.preRun(ctx, args); err != nil {
			return err
		}
	}
	return nil
}

func (command *Command) executePostRun(ctx context.Context, args []string) error {
	if command.postRun != nil {
		if err := command.postRun(ctx, args); err != nil {
			return err
		}
	}
	path := command.pathToRoot()
	for _, cmd := range path {
		if cmd.persistentPostRun != nil {
			if err := cmd.persistentPostRun(ctx, args); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Command) pathToRoot() []*Command {
	var command Command = (Command)(*c)
	path := []*Command{(*Command)(c)}
	for {
		if command.parent == nil {
			break
		}
		path = append(path, command.parent)
		command = *(command.parent)
	}
	return path
}
