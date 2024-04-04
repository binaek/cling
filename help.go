package cling

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func (c *Command) printHelp() error {
	parents := c.pathToRoot()
	slices.Reverse(parents)
	var cli *CLI
	parentsStr := make([]string, len(parents))
	for i, parent := range parents {
		parentsStr[i] = parent.name
		if parent.cli != nil {
			cli = parent.cli
		}
	}

	fmt.Println(c.longDescription)

	fmt.Println("Usage: ")
	usageString := fmt.Sprintf("%s %s", cli.name, strings.Join(parentsStr, " "))
	if len(c.children) > 0 {
		usageString = fmt.Sprintf("%s <command>", usageString)
	}

	if c.configInstance == nil {
		return nil
	}
	flags, args, err := extractConfigDefinitions(c.configInstance)
	if err != nil {
		return err
	}
	if len(args) > 0 {
		for _, arg := range args {
			if arg.optional {
				usageString = fmt.Sprintf("%s [%s]", usageString, arg.name)
				continue
			}
			usageString = fmt.Sprintf("%s <%s>", usageString, arg.name)
		}
	}

	if len(flags) > 0 {
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
	if len(flags) > 0 {
		buff := bytes.NewBuffer(nil)
		buff.WriteString("Flags:\n")
		flagsTable := tablewriter.NewWriter(buff)
		flagsTable.SetBorder(false)
		flagsTable.SetColumnSeparator("")
		for _, flag := range flags {
			flagsTable.Append(
				[]string{
					fmt.Sprintf("--%s", flag.name),
					flag.getHelpType(),
					flag.description,
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
