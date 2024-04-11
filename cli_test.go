package cling

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

type Config struct {
	Positional1 string `cling-name:"positional1"`
	Positional2 string `cling-name:"positional2"`
	Positional3 string `cling-name:"positional3"`
	StringFlag1 string `cling-name:"stringflag1"`
	IntFlag1    int    `cling-name:"intflag1"`
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("Config{\n")
	sb.WriteString(fmt.Sprintf("  Positional1: %s\n", c.Positional1))
	sb.WriteString(fmt.Sprintf("  Positional2: %s\n", c.Positional2))
	sb.WriteString(fmt.Sprintf("  Positional3: %s\n", c.Positional3))
	sb.WriteString(fmt.Sprintf("  StringFlag1: %s\n", c.StringFlag1))
	sb.WriteString(fmt.Sprintf("  IntFlag1: %d\n", c.IntFlag1))
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
		AddCommand(
			NewCommand("subcmd1", action).
				WithArgument(
					NewIntCmdInput("positional1").AsArgument(),
				).
				WithArgument(
					NewIntCmdInput("positional2").AsArgument(),
				).
				WithArgument(
					NewIntCmdInput("positional3").AsArgument(),
				).
				WithFlag(
					NewStringCmdInput("stringflag1").
						WithDefault("default").
						AsFlag(),
				).
				WithFlag(
					NewIntCmdInput("intflag1").Required(),
				),
		)

	ctx := context.Background()
	err := cli.Run(ctx, []string{
		"test",
		"subcmd1",
		"pos1",
		"pos2",
		"pos3",
		// "--stringflag1", "stringflag1",
		"--intflag1", "10",
	})
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
				WithChildCommand(
					NewCommand("subcmd11", action).
						WithArgument(
							NewStringCmdInput("arg1").
								WithDefault("default").
								WithValidators(
									NewStringLengthValidator(1, 10),
									NewEnumValidator(
										"one",
										"two",
										"three",
									),
								).
								Required().
								AsArgument(),
						).
						WithFlag(
							NewIntCmdInput("intflag").
								WithDefault(10).
								WithValidators(
									NewIntRangeValidator(0, 100),
									NewEnumValidator(1, 2, 3),
								).
								AsFlag(),
						),
				),
		).
		AddCommand(NewCommand("subcmd2", action))

	ctx := context.Background()
	err := cli.Run(ctx, []string{"test", "subcmd1", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
