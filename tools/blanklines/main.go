package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
)

// Exit codes.
const (
	exitOK      = 0
	exitFailure = 1
	exitUsage   = 2
)

// filePerm is used only if a fixed file's mode can't be read.
const filePerm = 0o600

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

// run checks or fixes the files named by args and returns the exit code. It collects
// its output and writes it once at the end, so a failed write is caught in one place.
func run(args []string, stdout, stderr io.Writer) int {

	fs := flag.NewFlagSet("blanklines", flag.ContinueOnError)
	fs.SetOutput(stderr)
	write := fs.Bool("w", false, "add the missing blank lines instead of listing them")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	paths := fs.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	var out, errs strings.Builder
	code := runFiles(paths, *write, &out, &errs)

	if _, err := io.WriteString(stderr, errs.String()); err != nil {
		return exitFailure
	}

	if _, err := io.WriteString(stdout, out.String()); err != nil {
		return exitFailure
	}

	return code
}

// runFiles processes every Go file under paths and returns the exit code.
func runFiles(paths []string, write bool, out, errs *strings.Builder) int {

	files, err := goFiles(paths)
	if err != nil {
		errs.WriteString("blanklines: " + err.Error() + "\n")

		return exitFailure
	}

	missing, touched, failed := 0, 0, false
	for _, file := range files {
		n, err := process(file, write, out)
		if err != nil {
			errs.WriteString("blanklines: " + err.Error() + "\n")

			failed = true
		}

		if n > 0 {
			missing += n
			touched++
		}
	}

	switch {
	case failed:
		return exitFailure
	case missing == 0:
		fmt.Fprintf(out, "checked %s: no blank lines missing\n", count(len(files), "file"))

		return exitOK
	case write:
		fmt.Fprintf(out, "fixed %d of %s: added %s\n", touched, count(len(files), "file"), count(missing, "blank line"))

		return exitOK
	default:
		fmt.Fprintf(out, "checked %s: %s missing in %d\n", count(len(files), "file"), count(missing, "blank line"), touched)

		return exitFailure
	}
}

// process checks one file, and with write set fixes it. It returns how many blank lines
// were missing, and writes each finding (or one line per fixed file) to out.
func process(file string, write bool, out *strings.Builder) (int, error) {

	src, err := os.ReadFile(file) //nolint:gosec // G304: reading the files the user names is the point
	if err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	if generated(src) {
		return 0, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, parser.SkipObjectResolution)
	if err != nil {
		return 0, fmt.Errorf("parse: %w", err)
	}

	findings := check(fset, file, f, strings.Split(string(src), "\n"))
	if len(findings) == 0 {
		return 0, nil
	}

	if !write {
		for _, fd := range findings {
			out.WriteString(fd.String() + "\n")
		}

		return len(findings), nil
	}

	return len(findings), writeFixed(file, src, findings, out)
}

// writeFixed adds the blank lines to file, keeping its mode, and records what it did.
func writeFixed(file string, src []byte, findings []finding, out *strings.Builder) error {

	fixed, err := fix(src, findings)
	if err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}

	mode := os.FileMode(filePerm)
	if info, statErr := os.Stat(file); statErr == nil {
		mode = info.Mode().Perm()
	}

	//nolint:gosec // G703: rewriting the files the user names, in place, is what -w does
	if err = os.WriteFile(file, fixed, mode); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	fmt.Fprintf(out, "fixed %s: added %s\n", file, count(len(findings), "blank line"))

	return nil
}

// count formats n with noun, plural unless n is 1: "1 file", "3 files".
func count(n int, noun string) string {

	if n == 1 {
		return "1 " + noun
	}

	return fmt.Sprintf("%d %ss", n, noun)
}
