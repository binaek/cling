package cling

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

type Config struct {
	Positional1 string   `cling-name:"positional1"`
	Positional2 string   `cling-name:"positional2"`
	Positional3 string   `cling-name:"positional3"`
	StringFlag1 string   `cling-name:"stringflag1"`
	IntFlag1    int      `cling-name:"intflag1"`
	SliceFlag1  []int    `cling-name:"sliceflag1"`
	SliceFlag2  []string `cling-name:"sliceflag2"`
	SliceFlag3  []string `cling-name:"sliceflag3"`
	SliceFlag4  []int    `cling-name:"sliceflag4"`
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
		WithCommand(
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
						WithValidator(NewEnumValidator("default", "one", "two", "three")).
						AsFlag().
						FromEnv([]string{"STRINGFLAG1"}),
				).
				WithFlag(
					NewCmdSliceInput[int]("sliceflag1").
						WithDefault([]int{1, 2, 3}).
						AsFlag(),
				).
				WithFlag(
					NewCmdSliceInput[string]("sliceflag2").
						WithDefault([]string{"one", "two", "three"}).
						AsFlag(),
				).
				WithFlag(
					NewCmdSliceInput[string]("sliceflag3").
						WithDefault([]string{"one", "two", "three"}).
						AsFlag(),
				).
				WithFlag(
					NewCmdSliceInput[int]("sliceflag4").
						WithDefault([]int{1, 2, 3}).
						AsFlag(),
				).
				WithFlag(
					NewIntCmdInput("intflag1").Required().AsFlag(),
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
		"--sliceflag1", "4",
		"--sliceflag1", "5",
		"--sliceflag1", "6",
		"--sliceflag2", "four",
		"--sliceflag2", "five",
		"--sliceflag2", "six",
		"--sliceflag3", "seven,eight,nine",
		"--sliceflag4", "7,8,9",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVersion(t *testing.T) {
	cli := NewCLI("test", "0.0.1").
		WithCommand(NewCommand("subcmd1", action)).
		WithCommand(NewCommand("subcmd2", action))

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
		WithCommand(
			NewCommand("subcmd1", action).
				WithChildCommand(
					NewCommand("subcmd11", action).
						WithArgument(
							NewStringCmdInput("arg1").
								WithDefault("default").
								WithValidator(
									ComposeValidator(
										NewStringLengthValidator(1, 10),
										NewEnumValidator(
											"one",
											"two",
											"three",
										),
									),
								).
								Required().
								AsArgument(),
						).
						WithFlag(
							NewIntCmdInput("intflag").
								WithDefault(10).
								WithValidator(
									ComposeValidator(
										NewIntRangeValidator(0, 100),
										NewEnumValidator(1, 2, 3),
									),
								).
								AsFlag(),
						),
				),
		).
		WithCommand(NewCommand("subcmd2", action))

	ctx := context.Background()
	err := cli.Run(ctx, []string{"test", "subcmd1", "--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
