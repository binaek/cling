package cling

import (
	"io"
	"os"
)

type CLI struct {
	name            string
	description     string
	longDescription string
	version         string
	commands        []*Command

	preRun  CommandHook
	postRun CommandHook

	stdout io.Writer
	stderr io.Writer
}

// WithCommand adds a command to the CLI
func (cli *CLI) WithCommand(command *Command) *CLI {
	cli.commands = append(cli.commands, command)
	return cli
}

// WithPreRun sets the pre-run hook for the CLI
func (cli *CLI) WithPreRun(hook CommandHook) *CLI {
	cli.preRun = hook
	return cli
}

// WithPostRun sets the post-run hook for the CLI
func (cli *CLI) WithPostRun(hook CommandHook) *CLI {
	cli.postRun = hook
	return cli
}

// WithDescription sets the description of the CLI
func (cli *CLI) WithDescription(description string) *CLI {
	cli.description = description
	return cli
}

// WithLongDescription sets the long description of the CLI
func (cli *CLI) WithLongDescription(longDescription string) *CLI {
	cli.longDescription = longDescription
	return cli
}

// NewCLI creates a new CLI
func NewCLI(name string, version string) *CLI {
	cli := &CLI{
		name:    name,
		version: version,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		preRun:  NoOpHook,
		postRun: NoOpHook,
	}
	return cli
}
