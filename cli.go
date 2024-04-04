package cling

type CLI struct {
	name            string
	description     string
	longDescription string
	version         string
	commands        []*Command
}

// Assuming AddCommand is a method of a CLI struct, make it a generic method
func (cli *CLI) AddCommand(command *Command) *CLI {
	// Add the command to your CLI structure
	// This might involve casting to *Command[any] if necessary, but be careful with type safety
	command.cli = cli
	cli.commands = append(cli.commands, command)
	return cli
}

func (cli *CLI) WithDescription(description string) *CLI {
	cli.description = description
	return cli
}

func (cli *CLI) WithLongDescription(longDescription string) *CLI {
	cli.longDescription = longDescription
	return cli
}

func NewCLI(name string, version string) *CLI {
	cli := &CLI{
		name:    name,
		version: version,
	}
	return cli
}
