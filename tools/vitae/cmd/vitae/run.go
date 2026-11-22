package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/lint"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// run executes one vitae command. It never exits; main maps its error to an exit code.
func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errorUsage("no command")
	}
	cmd, args := args[0], args[1:]
	switch cmd {
	case "version":
		_, err := fmt.Fprintln(stdout, version())
		return err //nolint:wrapcheck // a failed write to stdout needs no context
	case "help", "-h", "-help", "--help":
		_, err := fmt.Fprint(stdout, usage)
		return err //nolint:wrapcheck // a failed write to stdout needs no context
	case "validate", "render", "list":
	default:
		return errorUsage("unknown command %q", cmd)
	}

	var o options
	fs := flag.NewFlagSet("vitae "+cmd, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&o.input, "i", "resume/resume.yaml", "the résumé")
	fs.StringVar(&o.private, "private", "", "private overlay to merge (phone); local builds only")
	fs.BoolVar(&o.strict, "strict", false, "treat warnings as errors")
	fs.StringVar(&o.variant, "variant", "all", "variant to render, or all")
	fs.StringVar(&o.format, "format", "all", "format to render, or all")
	fs.StringVar(&o.outDir, "o", ".", "output directory, or - for standard output")
	if err := fs.Parse(args); err != nil {
		return errorUsage("%v", err)
	}

	r, err := resume.Load(o.input)
	if err != nil {
		return fmt.Errorf("load résumé: %w", err)
	}
	switch cmd {
	case "list":
		return list(r, fs.Args(), stdout)
	case "validate":
		return validate(r, &o, stdout, stderr)
	default:
		return renderAll(r, &o, stdout, stderr)
	}
}

// check prints every lint issue and fails if they include errors (or, when strict,
// warnings). Then it merges the private overlay, if one was given.
func check(r *resume.Resume, o *options, stderr io.Writer) error {
	issues := lint.Check(r)
	var report strings.Builder
	for _, i := range issues {
		report.WriteString(i.String() + "\n")
	}
	if _, err := io.WriteString(stderr, report.String()); err != nil {
		return fmt.Errorf("report issues: %w", err)
	}
	if lint.HasErrors(issues, o.strict) {
		return &errorExit{code: exitFailure, msg: fmt.Sprintf("%s: %d issues found", o.input, len(issues))}
	}
	if o.private != "" {
		if err := resume.ApplyPrivate(r, o.private); err != nil {
			return fmt.Errorf("load private overlay: %w", err)
		}
	}
	return nil
}
