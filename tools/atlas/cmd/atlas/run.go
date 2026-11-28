package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/layout"
	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/plate"
	"github.com/jonreyna/jonreyna.dev/tools/atlas/internal/render"
)

// Exit codes (see doc.go).
const (
	exitFailure = 1
	exitUsage   = 2
)

// The modes of the files and directories written: public web pages, readable by all.
const (
	filePerm = 0o644
	dirPerm  = 0o755
)

// Errors run returns.
var (
	// errUsage marks errors in the command line, which exit with exitUsage.
	errUsage = errors.New("usage")
	// errPlate marks a plate with problems, which run has listed on stderr.
	errPlate = errors.New("the plate has problems")
	// errCrossings marks a layout over the -max-crossings budget.
	errCrossings = errors.New("too many crossings")
)

// run dispatches to the validate or build command.
func run(args []string, stdout, stderr io.Writer) error {

	if len(args) == 0 {
		return fmt.Errorf("%w: atlas validate|build [flags] PLATE.yaml", errUsage)
	}

	switch args[0] {
	case "validate":
		return validate(args[1:], stdout, stderr)
	case "build":
		return build(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("%w: unknown command %q", errUsage, args[0])
	}
}

// validate reads a plate and reports its problems or a summary.
func validate(args []string, stdout, stderr io.Writer) error {

	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	requireChecked := fs.Bool("require-checked", false, "report links that aren't checked")

	path, err := parseArgs(fs, args)
	if err != nil {
		return err
	}

	p, err := load(path, *requireChecked, stderr)
	if err != nil {
		return err
	}

	unchecked := 0
	for i := range p.Links {
		if p.Links[i].Checked == nil {
			unchecked++
		}
	}

	return report(stdout, "validated %s: %d people, %d threads, %d links (%d unchecked)\n",
		filepath.Base(path), len(p.People), len(p.Threads), len(p.Links), unchecked)
}

// build validates and lays out a plate and writes its markup and stylesheet.
func build(args []string, stdout, stderr io.Writer) error {

	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(stderr)
	requireChecked := fs.Bool("require-checked", false, "report links that aren't checked")
	maxCrossings := fs.Int("max-crossings", -1, "fail above this many crossings; -1 for no limit")
	out := fs.String("o", "", "the markup `file` to write")
	css := fs.String("css", "", "the stylesheet `file` to write")

	path, err := parseArgs(fs, args)
	if err != nil {
		return err
	}

	if *out == "" || *css == "" {
		return fmt.Errorf("%w: -o and -css are required", errUsage)
	}

	p, err := load(path, *requireChecked, stderr)
	if err != nil {
		return err
	}

	o := layout.Defaults()

	l, err := layout.New(p, &o)
	if err != nil {
		return fmt.Errorf("lay out %s: %w", path, err)
	}

	if *maxCrossings >= 0 && l.Crossings > *maxCrossings {
		return fmt.Errorf("lay out %s: %w: %d, more than the %d allowed", path, errCrossings, l.Crossings, *maxCrossings)
	}

	var markup, rules bytes.Buffer

	err = errors.Join(render.HTML(&markup, p, l), render.CSS(&rules, p, l))
	if err != nil {
		return fmt.Errorf("render %s: %w", path, err)
	}

	for _, f := range []struct {
		path string
		data []byte
	}{{*out, markup.Bytes()}, {*css, rules.Bytes()}} {
		err = write(f.path, f.data, stdout)
		if err != nil {
			return err
		}
	}

	return report(stdout, "laid out %s: %d crossings\n", filepath.Base(path), l.Crossings)
}

// parseArgs parses flags and returns the one plate path that must follow them.
func parseArgs(fs *flag.FlagSet, args []string) (string, error) {

	err := fs.Parse(args)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errUsage, err)
	}

	if fs.NArg() != 1 {
		return "", fmt.Errorf("%w: give one plate file", errUsage)
	}

	return fs.Arg(0), nil
}

// load reads, parses and validates a plate, listing its problems on stderr.
func load(path string, requireChecked bool, stderr io.Writer) (*plate.Plate, error) {

	data, err := os.ReadFile(path) //nolint:gosec // G304: reading the plate the user names is the point
	if err != nil {
		return nil, fmt.Errorf("read plate: %w", err)
	}

	p, err := plate.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	issues := plate.Validate(p, requireChecked)
	for _, i := range issues {
		_, err = fmt.Fprintf(stderr, "%s: %s\n", filepath.Base(path), i)
		if err != nil {
			return nil, fmt.Errorf("report: %w", err)
		}
	}

	if len(issues) > 0 {
		return nil, fmt.Errorf("%s: %w (%d)", path, errPlate, len(issues))
	}

	return p, nil
}

// write writes a file, making its directory, and reports it.
func write(path string, data []byte, stdout io.Writer) error {

	err := os.MkdirAll(filepath.Dir(path), dirPerm)
	if err != nil {
		return fmt.Errorf("make directory for %s: %w", path, err)
	}

	err = os.WriteFile(path, data, filePerm)
	if err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return report(stdout, "wrote %s (%d bytes)\n", path, len(data))
}

// report prints a status line.
func report(stdout io.Writer, format string, args ...any) error {

	_, err := fmt.Fprintf(stdout, format, args...)
	if err != nil {
		return fmt.Errorf("report: %w", err)
	}

	return nil
}
