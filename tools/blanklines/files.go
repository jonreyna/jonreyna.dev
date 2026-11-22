package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// goFiles expands paths into the Go files to check, in walk order, skipping testdata,
// vendor and hidden directories.
func goFiles(paths []string) ([]string, error) {

	var out []string
	for _, root := range paths {
		info, err := os.Stat(root)
		if err != nil {
			return nil, fmt.Errorf("find files: %w", err)
		}

		if !info.IsDir() {
			out = append(out, root)

			continue
		}

		err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {

			if err != nil {
				return err
			}

			if d.IsDir() && path != root && skipDir(d.Name()) {
				return filepath.SkipDir
			}

			if !d.IsDir() && strings.HasSuffix(path, ".go") {
				out = append(out, path)
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("walk %s: %w", root, err)
		}
	}

	return out, nil
}

// skipDir reports whether a directory is never checked.
func skipDir(name string) bool {
	return name == "testdata" || name == "vendor" || strings.HasPrefix(name, ".")
}

// generated reports whether src carries the standard "Code generated ... DO NOT EDIT."
// marker.
func generated(src []byte) bool {

	for line := range bytes.Lines(src) {
		if bytes.HasPrefix(line, []byte("// Code generated ")) && bytes.Contains(line, []byte("DO NOT EDIT.")) {
			return true
		}

		if bytes.HasPrefix(line, []byte("package ")) {
			return false
		}
	}

	return false
}
