package cling

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var ErrInvalidCLIngConfig = errors.New("invalid CLIng configuration")

// Run executes the CLI with the given command line arguments.
func (c *CLI) Run(ctx context.Context, args []string) error {
	// get the executable name
	exec, err := os.Executable()
	if err != nil {
		return errors.Wrap(ErrInvalidCLIngConfig, "could not resolve executable name")
	}
	c.name = filepath.Base(exec)

	// validate the CLI
	if err := c.validate(); err != nil {
		return err
	}

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
		return command.printHelp(c)
	}

	// verify that there are no required arguments after an optional one
	for idx, arg := range command.arguments {
		isThisRequired := arg.isRequired()
		isThereAnymore := idx < (len(command.arguments) - 1)
		if !isThereAnymore {
			break
		}
		isNextRequired := command.arguments[idx+1].isRequired()

		if !isThisRequired && isThereAnymore && isNextRequired {
			return fmt.Errorf("required argument %s after optional argument %s", arg.Name(), command.arguments[idx+1].Name())
		}
	}

	// verify that all flags are either required or have a default value
	for _, flag := range command.flags {
		if !flag.isRequired() && !flag.hasDefault() {
			return fmt.Errorf("flag %s has no default value but is not required", flag.Name())
		}
	}

	newArgs := append(positionals, reconstructCmdLineFromFlags(flags)...)

	ctx = contextWithCommand(ctx, command)
	if c.preRun != nil {
		if err := c.preRun(ctx, newArgs); err != nil {
			return err
		}
	}
	err = command.execute(ctx, newArgs)
	if err != nil {
		if errors.Is(err, ErrInvalidCommand) {
			c.printUsage()
		}
		return err
	}
	if c.postRun != nil {
		if err := c.postRun(ctx, newArgs); err != nil {
			return err
		}
	}
	return nil
}

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

func (c *CLI) validate() error {
	if c.stdout == nil {
		return errors.Wrap(ErrInvalidCLIngConfig, "stdout channel not set")
	}

	if c.stderr == nil {
		return errors.Wrap(ErrInvalidCLIngConfig, "stderr channel not set")
	}

	for _, cmd := range c.commands {
		if err := cmd.validate(); err != nil {
			return err
		}
	}
	return nil
}
