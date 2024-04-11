package cling

import (
	"context"
)

type CommandHandler func(ctx context.Context, args []string) error

type Command struct {
	name            string
	description     string
	longDescription string
	shortcut        string
	action          CommandHandler
	flags           map[string]CmdFlag
	arguments       []CmdArg
	children        []*Command
	parent          *Command
	cli             *CLI
}

func (c *Command) WithDescription(description string) *Command {
	c.description = description
	return c
}

func (c *Command) WithLongDescription(longDescription string) *Command {
	c.longDescription = longDescription
	return c
}

func (command *Command) WithCommandShortcut(shortcut string) *Command {
	command.shortcut = shortcut
	return command
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

func NewCommand(name string, action CommandHandler) *Command {
	command := &Command{
		name:      name,
		action:    action,
		flags:     make(map[string]CmdFlag),
		arguments: []CmdArg{},
	}
	return command
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
