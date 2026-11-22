package main

import "fmt"

// Exit codes.
const (
	exitFailure = 1 // validation failed, or a file couldn't be read or written
	exitUsage   = 2 // bad command line
)

// exitError carries an exit code; an empty msg means the details were already printed.
type exitError struct {
	code int
	msg  string
}

func (e *exitError) Error() string { return e.msg }

// usageError returns an error that exits with exitUsage after printing the usage.
func usageError(format string, args ...any) error {
	return &exitError{code: exitUsage, msg: fmt.Sprintf(format, args...)}
}

const usage = `usage:
  vitae validate [-i FILE] [-private FILE] [-strict]
  vitae render   [-i FILE] [-private FILE] [-strict] [-variant NAME|all] [-format NAME|all] [-o DIR|-]
  vitae list     [-i FILE] variants|tags
  vitae version
`
