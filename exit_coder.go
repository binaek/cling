package cling

import "os"

type ExitCoder interface {
	ExitCode() int
}

type exitCoder struct {
	code int
	err  error
}

func (e *exitCoder) ExitCode() int {
	return e.code
}

func (e *exitCoder) Error() string {
	return e.err.Error()
}

func NewExitCoder(err error, code int) ExitCoder {
	return &exitCoder{
		code: code,
		err:  err,
	}
}

// ExitWithErrorMessage exits the program with a non-zero exit code if the given error is non-nil.
// If the given error is an `ExitCoder`, the exit code will be taken from the error, otherwise it will be 1.
// The error message will be printed to stderr as "Error: <message>\n".
//
// Uses `os.Exit` to exit the program. This function should be used only after when all cleanups are done.
func ExitWithMessage(err error) {
	exitWithError(err, true)
}

// ExitWithErrorCode - exits the program with a non-zero exit code if the given error is non-nil.
// If the given error is an `ExitCoder`, the exit code will be taken from the error, otherwise it will be 1.
//
// Uses `os.Exit` to exit the program. This function should be used only after when all cleanups are done.
func Exit(err error) {
	exitWithError(err, false)
}

func exitWithError(err error, printMessage bool) {
	if err == nil {
		return
	}
	if printMessage {
		_, _ = os.Stderr.WriteString("Error: " + err.Error() + "\n")
	}
	exitCode := 1
	if exitErr, ok := err.(ExitCoder); ok {
		exitCode = exitErr.ExitCode()
	}
	os.Exit(exitCode)
}
