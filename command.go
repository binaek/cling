package cling

import (
	"context"
	"slices"

	"github.com/pkg/errors"
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
	flags           []CmdFlag
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
		flags:             []CmdFlag{},
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
	command.flags = append(command.flags, flag)
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

var ErrInvalidCommand = errors.New("invalid command")

func (c *Command) validate() error {
	if c.name == "" {
		return errors.Wrapf(ErrInvalidCommand, "command name is required")
	}
	if c.action == nil {
		return errors.Wrapf(ErrInvalidCommand, "command action is required in command '%s'", c.name)
	}
	if err := c.validateFlagsAndArgs(); err != nil {
		return err
	}
	for _, child := range c.children {
		if err := child.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) validateFlagsAndArgs() error {
	// also validate that there's no conflict on the names across flags and arguments
	flagAndArgNames := make([]string, 0, len(c.flags))
	for _, flag := range c.flags {
		flagAndArgNames = append(flagAndArgNames, flag.Name())
	}
	for _, arg := range c.arguments {
		flagAndArgNames = append(flagAndArgNames, arg.Name())
	}
	slices.Sort(flagAndArgNames)
	cmpctNames := slices.Compact(flagAndArgNames)
	if len(flagAndArgNames) != len(cmpctNames) {
		return errors.Wrapf(ErrInvalidCommand, "duplicate flag and argument names found: %v", flagAndArgNames)
	}

	if err := c.validateFlags(); err != nil {
		return err
	}
	if err := c.validateArguments(); err != nil {
		return err
	}
	return nil
}

func (c *Command) validateFlags() error {
	names := make([]string, 0, len(c.flags))
	for _, flag := range c.flags {
		names = append(names, flag.Name())
	}
	slices.Sort(names)
	cmpctNames := slices.Compact(names)
	if len(names) != len(cmpctNames) {
		return errors.Wrapf(ErrInvalidCommand, "duplicate flag names found: %v", names)
	}
	return nil
}

func (c *Command) validateArguments() error {
	names := make([]string, 0, len(c.arguments))
	for _, arg := range c.arguments {
		names = append(names, arg.Name())
	}
	slices.Sort(names)
	cmpctNames := slices.Compact(names)
	if len(names) != len(cmpctNames) {
		return errors.Wrapf(ErrInvalidCommand, "duplicate argument names found: %v", names)
	}
	return nil
}
