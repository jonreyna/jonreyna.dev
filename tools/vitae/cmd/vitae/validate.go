package main

import (
	"fmt"
	"io"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// validate checks the résumé and reports what it validated.
func validate(r *resume.Resume, o *options, stdout, stderr io.Writer) error {

	if err := check(r, o, stderr); err != nil {
		return err
	}

	_, err := fmt.Fprintf(stdout, "validated %s: %d roles, %d variants\n", o.input, len(r.Experience), len(r.Variants))
	return err //nolint:wrapcheck // a failed write to stdout needs no context
}
