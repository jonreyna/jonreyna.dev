// Command vitae validates resume.yaml and renders it, one file per variant and format.
//
//	vitae validate [-i FILE] [--private FILE] [--strict]
//	vitae render   [-i FILE] [--private FILE] [--variant NAME|all] [--format json] [-o DIR|-]
//	vitae list     variants|tags
//	vitae version
//
// Exit codes: 0 success, 1 validation failure, 2 usage error.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/filter"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/lint"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/model"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/render"
)

const usage = `usage:
  vitae validate [-i FILE] [--private FILE] [--strict]
  vitae render   [-i FILE] [--private FILE] [--variant NAME|all] [--format json|all] [-o DIR|-]
  vitae list     [-i FILE] variants|tags
  vitae version
`

// Formats vitae can render. LaTeX, text and Markdown come next (Phase 4, steps 3–4).
var formats = []string{"json"}

type exitError struct {
	code int
	msg  string
}

func (e *exitError) Error() string { return e.msg }

func usageErr(format string, args ...any) error {
	return &exitError{2, fmt.Sprintf(format, args...)}
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		code := 1
		var ee *exitError
		if errors.As(err, &ee) {
			code = ee.code
			if code == 2 {
				fmt.Fprint(os.Stderr, usage)
			}
		}
		if err.Error() != "" {
			fmt.Fprintln(os.Stderr, "vitae:", err)
		}
		os.Exit(code)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return usageErr("no command")
	}
	cmd, args := args[0], args[1:]
	fs := flag.NewFlagSet("vitae "+cmd, flag.ContinueOnError)
	fs.SetOutput(stderr)
	input := fs.String("i", "resume/resume.yaml", "the résumé")
	private := fs.String("private", "", "private overlay to merge (phone); local builds only")
	strict := fs.Bool("strict", false, "treat warnings as errors")
	variant := fs.String("variant", "all", "variant to render, or all")
	format := fs.String("format", "all", "format to render, or all")
	outDir := fs.String("o", ".", "output directory, or - for standard output")

	switch cmd {
	case "version":
		fmt.Fprintln(stdout, version())
		return nil
	case "validate", "render", "list":
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return nil
	default:
		return usageErr("unknown command %q", cmd)
	}
	if err := fs.Parse(args); err != nil {
		return usageErr("%v", err)
	}

	r, err := model.Load(*input)
	if err != nil {
		return err
	}
	issues := lint.Check(r)

	switch cmd {
	case "list":
		return list(r, fs.Args(), stdout)
	case "validate":
		for _, i := range issues {
			fmt.Fprintln(stderr, i)
		}
		if lint.HasErrors(issues, *strict) {
			return &exitError{1, ""}
		}
		if *private != "" {
			if err := model.ApplyPrivate(r, *private); err != nil {
				return err
			}
		}
		fmt.Fprintf(stdout, "%s: OK (%d roles, %d variants)\n", *input, len(r.Experience), len(r.Variants))
		return nil
	}

	// render: refuse to render a résumé with errors; warnings are printed and allowed.
	for _, i := range issues {
		fmt.Fprintln(stderr, i)
	}
	if lint.HasErrors(issues, *strict) {
		return &exitError{1, "not rendering: fix the errors above"}
	}
	if *private != "" {
		if err := model.ApplyPrivate(r, *private); err != nil {
			return err
		}
	}
	variants, err := pick(*variant, sortedKeys(r.Variants), "variant")
	if err != nil {
		return err
	}
	fmts, err := pick(*format, formats, "format")
	if err != nil {
		return err
	}
	if *outDir == "-" && len(variants)*len(fmts) != 1 {
		return usageErr("-o - writes one file; choose one --variant and one --format")
	}
	for _, name := range variants {
		resolved, err := filter.Apply(r, name)
		if err != nil {
			return err
		}
		for _, f := range fmts {
			data, err := renderOne(resolved, f)
			if err != nil {
				return err
			}
			if *outDir == "-" {
				_, err = stdout.Write(data)
				return err
			}
			if err := os.MkdirAll(*outDir, 0o755); err != nil {
				return err
			}
			path := filepath.Join(*outDir, fmt.Sprintf("resume.%s.%s", name, f))
			if err := os.WriteFile(path, data, 0o644); err != nil {
				return err
			}
			fmt.Fprintln(stdout, path)
		}
	}
	return nil
}

func renderOne(r *filter.Resolved, format string) ([]byte, error) {
	switch format {
	case "json":
		return render.JSON(r)
	}
	return nil, fmt.Errorf("format %q is not implemented", format)
}

func pick(choice string, all []string, what string) ([]string, error) {
	if choice == "all" {
		return all, nil
	}
	if !slices.Contains(all, choice) {
		return nil, usageErr("unknown %s %q (have: %s)", what, choice, strings.Join(all, ", "))
	}
	return []string{choice}, nil
}

func list(r *model.Resume, args []string, w io.Writer) error {
	if len(args) != 1 {
		return usageErr("list needs one of: variants, tags")
	}
	switch args[0] {
	case "variants":
		for _, name := range sortedKeys(r.Variants) {
			fmt.Fprintf(w, "%-10s %s\n", name, r.Variants[name].Description)
		}
	case "tags":
		counts := map[string]int{}
		for _, role := range r.Experience {
			for _, h := range role.Highlights {
				for _, t := range h.Tags {
					counts[t]++
				}
			}
		}
		for _, t := range sortedKeys(r.Tags) {
			fmt.Fprintf(w, "%-14s %2d  %s\n", t, counts[t], r.Tags[t])
		}
	default:
		return usageErr("list needs one of: variants, tags")
	}
	return nil
}

func version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				return "vitae " + s.Value[:7]
			}
		}
	}
	return "vitae dev"
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
