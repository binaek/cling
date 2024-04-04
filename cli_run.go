package cling

import (
	"context"
	"errors"
	"fmt"
)

func (c *CLI) Run(ctx context.Context, args []string) error {
	// parse the arguments
	flags, positionals := parseArguments(args)
	if len(positionals) < 2 {
		c.printUsage()
		return errors.New("missing command")
	}

	// remove the executable name
	positionals = positionals[1:]

	// do we have a --version in the flags
	if _, ok := flags["version"]; ok {
		fmt.Printf("%s v%s\n", c.name, c.version)
		return nil
	}

	command := c.findCommand(positionals)
	if command == nil {
		c.printUsage()
		return nil
	}

	pathToRoot := command.pathToRoot()
	positionals = positionals[len(pathToRoot):]

	if _, ok := flags["help"]; ok {
		return command.printHelp()
	}

	newArgs := append(positionals, reconstructCmdLineFromFlags(flags)...)
	return command.action(ctx, newArgs)
}

func (c *CLI) printUsage() {}

func reconstructCmdLineFromFlags(f map[string][]string) []string {
	flags := []string{}
	for k, v := range f {
		for _, val := range v {
			flags = append(flags, fmt.Sprintf("--%s=%s", k, val))
		}
	}
	return flags
}

func (c *CLI) findCommand(names []string) *Command {
	for _, cmd := range c.commands {
		if cmd.name == names[0] {
			return c.findCmd(cmd, names, 1)
		}
	}
	return nil
}

func (c *CLI) findCmd(currentCMD *Command, names []string, index int) *Command {
	if index == len(names) {
		return currentCMD
	}
	for _, child := range currentCMD.children {
		if child.name == names[index] {
			return c.findCmd(child, names, index+1)
		}
	}
	return currentCMD
}
