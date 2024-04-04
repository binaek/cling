package cling

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

type Config struct {
	WorkspacePath     string   `name:"workspace-path1" position:"0" description:"The workspace path1"`
	WorkspacePath1    string   `name:"workspace-path2" position:"1" description:"The workspace path2"`
	WorkspacePath2    string   `name:"workspace-path3" position:"2" description:"The workspace path3"`
	WorkspaceDatabase string   `name:"workspace-database" short:"w" description:"The workspace database"`
	BoolVal           *bool    `name:"boolval" short:"b" description:"A boolean value"`
	BoolVal2          bool     `name:"boolval2" short:"B" description:"A boolean value"`
	StrSliceVal       []string `name:"strsliceval" short:"s" description:"A string slice value"`
	IntSliceVal       []int    `name:"intsliceval" short:"i" description:"An int slice value"`
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("Config{\n")
	sb.WriteString(fmt.Sprintf("  WorkspacePath: %s,\n", c.WorkspacePath))
	sb.WriteString(fmt.Sprintf("  WorkspacePath1: %s,\n", c.WorkspacePath1))
	sb.WriteString(fmt.Sprintf("  WorkspacePath2: %s,\n", c.WorkspacePath2))
	sb.WriteString(fmt.Sprintf("  WorkspaceDatabase: %s,\n", c.WorkspaceDatabase))
	sb.WriteString(fmt.Sprintf("  BoolVal: %t,\n", *c.BoolVal))
	sb.WriteString(fmt.Sprintf("  BoolVal2: %t,\n", c.BoolVal2))
	sb.WriteString(fmt.Sprintf("  StrSliceVal: %v,\n", c.StrSliceVal))
	sb.WriteString(fmt.Sprintf("  IntSliceVal: %v,\n", c.IntSliceVal))
	sb.WriteString("}")

	return sb.String()
}

func action(ctx context.Context, args []string) error {
	cfg := &Config{}
	if err := Hydrate(ctx, args, cfg); err != nil {
		return err
	}
	fmt.Println(cfg)
	return nil
}

func TestRun(t *testing.T) {
	cli := NewCLI("test", "0.0.1").
		AddCommand(NewCommand("subcmd1", action)).
		AddCommand(NewCommand("subcmd2", action))

	ctx := context.Background()
	err := cli.Run(ctx, []string{"test", "subcmd1", "pos1", "pos2", "pos3", "--strsliceval", "val1", "--strsliceval", "val2", "--intsliceval", "1", "--intsliceval", "2", "--workspace-database", "testdb", "--boolval", "--boolval2", "false", "workspace"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVersion(t *testing.T) {
	cli := NewCLI("test", "0.0.1").
		AddCommand(NewCommand("subcmd1", action)).
		AddCommand(NewCommand("subcmd2", action))

	ctx := context.Background()
	err := cli.Run(ctx, []string{"test", "subcmd1", "--version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHelp(t *testing.T) {
	cli := NewCLI("test", "0.0.1").
		WithDescription("Test CLI").
		WithLongDescription("This is a test CLI").
		AddCommand(
			NewCommand("subcmd1", action).
				WithHelpFrom(&Config{}).
				WithSubcommand(
					NewCommand("subcmd11", action).
						WithHelpFrom(&Config{}),
				),
		).
		AddCommand(NewCommand("subcmd2", action).WithHelpFrom(&Config{}))

	ctx := context.Background()
	err := cli.Run(ctx, []string{"test", "subcmd1", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
