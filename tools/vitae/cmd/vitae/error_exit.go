package main

import "fmt"

// Exit codes carried by errorExit.
const (
	exitFailure = 1 // validation failed, or a file couldn't be read or written
	exitUsage   = 2 // bad command line; main prints the usage
)

// errorExit is an error that sets the exit code. An empty msg means the details were
// already printed.
type errorExit struct {
	code int
	msg  string
}

func (e *errorExit) Error() string { return e.msg }

// errorUsage returns the usage version of errorExit (exit code exitUsage), so main prints the usage after msg.
func errorUsage(format string, args ...any) error {
	return &errorExit{code: exitUsage, msg: fmt.Sprintf(format, args...)}
}
