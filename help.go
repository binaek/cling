package cling

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func (c *Command) printHelp(ctx context.Context, cli *CLI) error {
	parents := c.pathToRoot()
	slices.Reverse(parents)
	parentsStr := make([]string, len(parents))
	for i, parent := range parents {
		parentsStr[i] = parent.name
	}

	fmt.Println(c.longDescription)

	fmt.Println("Usage: ")
	usageString := fmt.Sprintf("%s %s", cli.name, strings.Join(parentsStr, " "))
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
	fmt.Printf("  %s\n", usageString)

	// Print available commands if any
	if len(c.children) > 0 {
		buff := bytes.NewBuffer(nil)
		buff.WriteString("Available Commands:\n")
		for _, child := range c.children {
			buff.WriteString(
				fmt.Sprintf("  %s\t%s\n", child.name, child.description),
			)
		}
		fmt.Println()
		fmt.Print(buff.String())
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

		fmt.Println()
		fmt.Print(buff.String())
		fmt.Println()
	}

	if len(c.children) > 0 {
		fmt.Printf("Use \"%s %s [command] --help\" for more information about a command.\n", cli.name, strings.Join(parentsStr, " "))
	}

	return nil
}
