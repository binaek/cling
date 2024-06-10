package cling

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func (c *CLI) printUsage() {
	fmt.Fprintf(c.stdout, "Usage: %s [command] [flags] [arguments]\n\n", c.name)
	fmt.Fprintln(c.stdout, "Available Commands:")
	for _, cmd := range c.commands {
		fmt.Fprintf(c.stdout, "  %s\t%s\n", cmd.name, cmd.description)
		for _, child := range cmd.children { // Print subcommands
			fmt.Fprintf(c.stdout, "    %s\t%s\n", child.name, child.description)
		}
	}

	fmt.Fprintln(c.stdout, "\nFlags:")
	fmt.Fprintln(c.stdout, "  --help\tShow help information")
	fmt.Fprintln(c.stdout, "  --version\tShow version information")

	fmt.Fprintf(c.stdout, "\nUse \"%s [command] --help\" for more information about a command.\n", c.name)
}

func (c *Command) printHelp(cli *CLI) error {
	path2Root := c.pathToRoot()
	slices.Reverse(path2Root)
	pathStr := make([]string, len(path2Root))
	for i, parent := range path2Root {
		pathStr[i] = parent.name
	}

	fmt.Fprintln(cli.stdout, c.longDescription)

	fmt.Fprintln(cli.stdout, "Usage: ")
	usageString := fmt.Sprintf("%s %s", cli.name, strings.Join(pathStr, " "))
	if len(c.children) > 0 {
		usageString = fmt.Sprintf("%s <command>", usageString)
	}

	if len(c.arguments) > 0 {
		for _, arg := range c.arguments {
			if !arg.isRequired() {
				usageString = fmt.Sprintf("%s [%s]", usageString, arg.Name())
				continue
			}
			usageString = fmt.Sprintf("%s <%s>", usageString, arg.Name())
		}
	}

	if len(c.flags) > 0 {
		usageString = fmt.Sprintf("%s [flags]", usageString)
	}
	fmt.Fprintf(cli.stdout, "  %s\n", usageString)

	// Print available commands if any
	if len(c.children) > 0 {
		buff := bytes.NewBuffer(nil)
		buff.WriteString("Available Commands:\n")
		for _, child := range c.children {
			buff.WriteString(
				fmt.Sprintf("  %s\t%s\n", child.name, child.description),
			)
		}
		fmt.Fprintln(cli.stdout)
		fmt.Fprintln(cli.stdout, buff.String())
	}

	// Print flags
	if len(c.flags) > 0 {
		buff := bytes.NewBuffer(nil)
		buff.WriteString("Flags:\n")
		flagsTable := tablewriter.NewWriter(buff)
		flagsTable.SetBorder(false)
		flagsTable.SetColumnSeparator("")
		for _, flag := range c.flags {
			flagsTable.Append(
				[]string{
					fmt.Sprintf("--%s", flag.Name()),
					flag.Description(),
				},
			)
		}
		flagsTable.Render()

		fmt.Fprintln(cli.stdout)
		fmt.Fprintln(cli.stdout, buff.String())
		fmt.Fprintln(cli.stdout)
	}

	if len(c.children) > 0 {
		fmt.Fprintf(cli.stdout, "Use \"%s %s [command] --help\" for more information about a command.\n", cli.name, strings.Join(pathStr, " "))
	}

	return nil
}
