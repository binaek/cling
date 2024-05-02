package cling

type CLI struct {
	name            string
	description     string
	longDescription string
	version         string
	commands        []*Command
}

// Assuming WithCommand is a method of a CLI struct, make it a generic method
func (cli *CLI) WithCommand(command *Command) *CLI {
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
