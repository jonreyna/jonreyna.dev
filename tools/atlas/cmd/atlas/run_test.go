package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const plateFile = "../../testdata/minimal.yaml"

func TestRun(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name    string
		args    func(dir string) []string
		wantErr error
		wantOut string
	}{
		{name: "no command", args: func(string) []string { return nil }, wantErr: errUsage},
		{name: "unknown command", args: func(string) []string { return []string{"draw"} }, wantErr: errUsage},
		{name: "validate", args: func(string) []string { return []string{"validate", plateFile} },
			wantOut: "validated minimal.yaml: 4 people, 2 threads, 4 links (1 unchecked)"},
		{name: "validate requiring checks", args: func(string) []string {
			return []string{"validate", "-require-checked", plateFile}
		}, wantErr: errPlate},
		{name: "build without outputs", args: func(string) []string { return []string{"build", plateFile} },
			wantErr: errUsage},
		{name: "build", args: func(dir string) []string {
			return []string{"build", "-o", filepath.Join(dir, "a", "m.html"), "-css", filepath.Join(dir, "a", "m.css"), plateFile}
		}, wantOut: "laid out minimal.yaml: 0 crossings"},
		{name: "build a missing plate", args: func(dir string) []string {
			return []string{"build", "-o", filepath.Join(dir, "m.html"), "-css", filepath.Join(dir, "m.css"), "missing.yaml"}
		}, wantErr: os.ErrNotExist},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			var stdout, stderr bytes.Buffer

			err := run(tc.args(t.TempDir()), &stdout, &stderr)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("run() error = %v, want %v", err, tc.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("run() error = %v; stderr: %s", err, stderr.String())
			}

			if !strings.Contains(stdout.String(), tc.wantOut) {
				t.Fatalf("stdout = %q, want it to contain %q", stdout.String(), tc.wantOut)
			}
		})
	}
}
