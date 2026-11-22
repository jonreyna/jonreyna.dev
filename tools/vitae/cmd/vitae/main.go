package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {

	err := run(os.Args[1:], os.Stdout, os.Stderr)
	if err == nil {
		return
	}

	code := exitFailure
	if ee, ok := errors.AsType[*errorExit](err); ok {
		code = ee.code
	}

	if code == exitUsage {
		fmt.Fprint(os.Stderr, usage)
	}

	if msg := err.Error(); msg != "" {
		fmt.Fprintln(os.Stderr, "vitae:", msg)
	}

	os.Exit(code)
}
