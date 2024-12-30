# CLIng is a command-line parser for Go

The aim of CLIng is to support arbitrarily complex command line flags and arguments for Go Programs.

Command line inputs are expressed as Go annotated `structs` that are hydrated at runtime using the generic `cling.Hydrate` function.

## Usage

```go
package main

import (
  "context"
  "fmt"
  "os"
  "os/signal"
  "syscall"

  "github.com/binaek/cling"
  "github.com/pkg/errors"
)

// Version of the program
var version = "0.0.1-dev.0"

func main() {
  // create a context so that we can shut things down gracefully when we receive a SIGNAL
  // we can pass this context to the CLI and it will be used for all downstream calls
  ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL)

  cli := setupCLI(ctx, version)
  if err := cli.Run(ctx, os.Args); err != nil {
    fmt.Println("Error:", err)
  }
}

type startCmdInput struct {
  AnyArbitraryName     string `cling-name:"string-input"`
  AnotherArbitraryName int    `cling-name:"int-input"`
}

// This is the handler for the `start` command.
// It will be called when the `start` command is executed in the CLI.
// The `args` parameter will contain the arguments passed to the command without the executable.
func startCmdHandler(ctx context.Context, args []string) error {
  input := startCmdInput{}

  // Hydrate the input struct with the arguments passed to the command
  // This will populate the input struct with the values of the flags passed to the command
  if err := cling.Hydrate(ctx, args, &input); err != nil {
    return errors.Wrap(err, "failed to hydrate input arguments")
  }

  fmt.Println("String input:", input.AnyArbitraryName)
  fmt.Println("Int input:", input.AnotherArbitraryName)

  return nil
}

func setupCLI(ctx context.Context, version string) *cling.CLI {
  cli := cling.NewCLI("program", version).
    WithDescription("This is an example go program").
    WithCommand(
      cling.NewCommand("start", startCmdHandler).
        WithDescription("Start the program").
        WithFlag(
          cling.
            NewStringCmdInput("string-input"). // MUST match with the a field in the input struct
            Required().                        // Make this flag required
            WithDescription("A string input"). // Give a description for help content generation
            AsFlag(),                          // This is important to tell cling that this is a flag
        ).
        WithFlag(
          cling.
            NewIntCmdInput("int-input").    // MUST match with the a field in the input struct
            WithDefault(1000).              // Set a default value - this is not required
            WithDescription("A INT input"). // Give a description for help content generation
            AsFlag(),                       // This is important to tell cling that this is a flag
        ),
    )
  return cli
}

```

## Features

### Version

CLIng provides a `--version` flag that can be used to show the version of the CLI.

```shell
$ program --version
program version 0.0.1-dev.0
```

```shell
$ program start --version
program version 0.0.1-dev.0
```

### Help

CLIng provides a `--help` flag that can be used to show help content for the CLI or it's commands. This is automatically generated based on the input struct and the flags that are defined.

#### Help for the CLI

When the `--help` flag is specified, the help content is generated for the CLI and it's commands.

```shell
$ program --help
Usage: program [command] [flags] [arguments]

Available Commands:
  start

Flags:
  --help        Show help information
  --version     Show version information

Use "program [command] --help" for more information about a command.
```

#### Help for a Command

When the `--help` flag is specified with a command, the help content is generated for the command with it's flags and arguments.

```shell
$ program start --help
Usage: program start [flags] [arguments]

Flags:
  --string-input string   A string input
  --int-input int         A INT input

```
