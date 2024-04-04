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
	configInstance  any
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

func (command *Command) WithHelpFrom(configInstance any) *Command {
	command.configInstance = configInstance
	return command
}

func (command *Command) WithCommandShortcut(shortcut string) *Command {
	command.shortcut = shortcut
	return command
}

func (command *Command) WithSubcommand(subcommand *Command) *Command {
	command.children = append(command.children, subcommand)
	subcommand.parent = (*Command)(command)
	return command
}

func NewCommand(name string, action CommandHandler) *Command {
	command := &Command{
		name:   name,
		action: action,
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
