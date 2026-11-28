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
	if errors.Is(err, errUsage) {
		code = exitUsage
	}

	fmt.Fprintln(os.Stderr, "atlas:", err)
	os.Exit(code)
}
