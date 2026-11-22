package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/render"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/variant"
)

// Output permissions: the files are build artifacts for the owner.
const (
	dirPerm  = 0o750
	filePerm = 0o600
)

// renderAll validates the résumé, then writes each chosen variant in each chosen format,
// printing "wrote PATH" after each file.
func renderAll(r *resume.Resume, o *options, stdout, stderr io.Writer) error {
	if err := check(r, o, stderr); err != nil {
		return err
	}
	variants, err := pick(o.variant, slices.Sorted(maps.Keys(r.Variants)), "variant")
	if err != nil {
		return err
	}
	formats, err := pick(o.format, render.Formats(), "format")
	if err != nil {
		return err
	}
	if o.outDir == "-" && len(variants)*len(formats) != 1 {
		return usageError("-o - writes one file; choose one -variant and one -format")
	}
	for _, name := range variants {
		resolved, err := variant.Apply(r, name)
		if err != nil {
			return fmt.Errorf("apply variant: %w", err)
		}
		for _, format := range formats {
			if err := renderOne(resolved, format, o.outDir, stdout); err != nil {
				return err
			}
		}
	}
	return nil
}

// renderOne renders one variant in one format to outDir, or to stdout if outDir is "-".
func renderOne(resolved *variant.Resolved, format, outDir string, stdout io.Writer) error {
	renderer, err := render.For(format)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	data, err := renderer.Render(resolved)
	if err != nil {
		return fmt.Errorf("render %s as %s: %w", resolved.Variant, format, err)
	}
	if outDir == "-" {
		_, err = stdout.Write(data)
		return err //nolint:wrapcheck // a failed write to stdout needs no context
	}
	if err = os.MkdirAll(outDir, dirPerm); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	path := filepath.Join(outDir, fmt.Sprintf("resume.%s.%s", resolved.Variant, format))
	if err = os.WriteFile(path, data, filePerm); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	_, err = fmt.Fprintln(stdout, "wrote", path)
	return err //nolint:wrapcheck // a failed write to stdout needs no context
}

// pick returns every name in all for "all", or the one name chosen.
func pick(choice string, all []string, what string) ([]string, error) {
	if choice == "all" {
		return all, nil
	}
	if !slices.Contains(all, choice) {
		return nil, usageError("unknown %s %q (have: %s)", what, choice, strings.Join(all, ", "))
	}
	return []string{choice}, nil
}
